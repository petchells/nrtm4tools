package pg

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/petchells/nrtm4client/internal/nrtm4/persist"
	"github.com/petchells/nrtm4client/internal/nrtm4/pg/db"
	pgpersist "github.com/petchells/nrtm4client/internal/nrtm4/pg/persist"
	"github.com/petchells/nrtm4client/internal/nrtm4/rpsl"
	"github.com/petchells/nrtm4client/internal/nrtm4/util"
)

// PostgresRepository implementation of the Repository interface
type PostgresRepository struct {
}

// Initialize implementation of the Repository interface
func (repo PostgresRepository) Initialize(dbURL string) error {
	return db.InitializeConnectionPool(dbURL)
}

// GetSources returns a list of all sources
func (repo PostgresRepository) GetSources() ([]persist.NRTMSource, error) {
	var sources []persist.NRTMSource
	var err error
	var pgsources []pgpersist.NRTMSource
	src := new(pgpersist.NRTMSource)
	err = db.WithTransaction(func(tx pgx.Tx) error {
		pgsources, err = db.GetAll(tx, *src, nil)
		return err
	})
	if err != nil {
		logger.Error("Error in GetSources", "error", err)
		return sources, err
	}
	for _, s := range pgsources {
		sources = append(sources, s.AsNRTMSource())
	}
	return sources, nil
}

// GetNotificationHistory gets the last 100 notification versions
func (repo PostgresRepository) GetNotificationHistory(source persist.NRTMSource, fromVersion, toVersion uint32) ([]persist.Notification, error) {
	if toVersion < fromVersion {
		return []persist.Notification{}, nil
	}
	notif := new(pgpersist.Notification)
	notifDesc := db.GetDescriptor(notif)
	sql := fmt.Sprintf(`
		SELECT %v
		FROM %v
		WHERE nrtm_source_id = $1
		AND version >= $2
		AND version <= $3
		ORDER BY version DESC
		LIMIT 1000
		`,
		notifDesc.ColumnNamesCommaSeparated(),
		notifDesc.TableName(),
	)
	notifs := make([]persist.Notification, 0, 100)
	err := db.WithTransaction(func(tx pgx.Tx) error {
		rows, err := tx.Query(context.Background(), sql, source.ID, fromVersion, toVersion)
		if err != nil {
			return err
		}
		for rows.Next() {
			ent := *notif
			err = rows.Scan(db.SelectValues(&ent)...)
			if err != nil {
				return err
			}
			notifs = append(notifs, asNotification(ent))
		}
		return nil
	})
	return notifs, err
}

// SaveSource updates a source if ID is non-zero, or creates a new one if it is
func (repo PostgresRepository) SaveSource(source persist.NRTMSource, notification persist.NotificationJSON) (persist.NRTMSource, error) {
	var pgSource pgpersist.NRTMSource
	err := db.WithTransaction(func(tx pgx.Tx) error {
		if source.ID == 0 {
			pgSource = pgpersist.NewNRTMSource(source)
			return db.Create(tx, &pgSource)
		}
		pgSource = pgpersist.FromNRTMSource(source)
		err := db.Update(tx, &pgSource)
		if err != nil {
			return err
		}
		return pgpersist.NewNotification(tx, source.ID, notification)
	})
	return pgSource.AsNRTMSource(), err
}

// Close implementation of the interface. Nothing needed for pg (for now)
func (repo PostgresRepository) Close() error {
	return nil
}

// SaveFile saves a reference to an NRTM file
func (repo PostgresRepository) SaveFile(nrtmFile *persist.NRTMFile) error {
	return db.WithTransaction(func(tx pgx.Tx) error {
		st := pgpersist.NRTMFile{
			ID:           uint64(db.NextID()),
			Version:      nrtmFile.Version,
			URL:          nrtmFile.URL,
			Type:         nrtmFile.Type.String(),
			NRTMSourceID: nrtmFile.NrtmSourceID,
			FileName:     nrtmFile.FileName,
			Created:      util.AppClock.Now(),
		}
		nrtmFile.ID = st.ID
		return db.Create(tx, &st)
	})
}

// SaveSnapshotObjects saves a list of rpsl object in a go routine
func (repo PostgresRepository) SaveSnapshotObjects(
	source persist.NRTMSource,
	rpslObjects []rpsl.Rpsl,
	file persist.NrtmFileJSON,
) error {
	if len(rpslObjects) == 0 {
		return nil
	}
	return db.WithTransaction(func(tx pgx.Tx) error {
		inputRows := make([][]any, len(rpslObjects))
		for i, rpslObject := range rpslObjects {
			inputRow := []any{
				uint64(db.NextID()),
				rpslObject.ObjectType,
				rpslObject.PrimaryKey,
				source.ID,
				file.Version,
				0,
				rpslObject.Payload,
			}
			inputRows[i] = inputRow
		}
		rpslDescriptor := db.GetDescriptor(&pgpersist.RPSLObject{})
		_, err := tx.CopyFrom(
			context.Background(),
			pgx.Identifier{rpslDescriptor.TableName()},
			rpslDescriptor.ColumnNames(),
			pgx.CopyFromRows(inputRows),
		)
		if err != nil {
			types := util.NewSet[string]()
			for _, inp := range rpslObjects {
				types.Add(inp.ObjectType)
			}
			logger.Warn("Failed to save objects", "types", types.String(), "error", err)
			return err
		}
		return nil
	})
}

// AddModifyObject updates an RPSL object by setting `to_version` and inserting a new row
func (repo PostgresRepository) AddModifyObject(
	source persist.NRTMSource,
	rpsl rpsl.Rpsl,
	file persist.NrtmFileJSON,
) error {
	rpslObject := new(pgpersist.RPSLObject)
	newRow := pgpersist.RPSLObject{
		ID:           db.NextID(),
		ObjectType:   rpsl.ObjectType,
		PrimaryKey:   rpsl.PrimaryKey,
		NRTMSourceID: source.ID,
		FromVersion:  file.Version,
		RPSL:         rpsl.Payload,
	}
	return db.WithTransaction(func(tx pgx.Tx) error {
		err := deleteObject(tx, source, rpslObject, rpsl.ObjectType, rpsl.PrimaryKey, file)
		if err == nil || err == pgx.ErrNoRows {
			return db.Create(tx, &newRow)
		}
		return err
	})
}

// DeleteObject doesn't remove any rows, instead it sets `to_version` to the file version
func (repo PostgresRepository) DeleteObject(
	source persist.NRTMSource,
	objectType string,
	primaryKey string,
	file persist.NrtmFileJSON,
) error {
	rpslObject := new(pgpersist.RPSLObject)
	return db.WithTransaction(func(tx pgx.Tx) error {
		return deleteObject(tx, source, rpslObject, objectType, primaryKey, file)
	})
}

func deleteObject(
	tx pgx.Tx,
	source persist.NRTMSource,
	rpslObject *pgpersist.RPSLObject,
	objectType string,
	primaryKey string,
	file persist.NrtmFileJSON,
) error {
	sql := selectObjectQuery()
	err := tx.QueryRow(context.Background(), sql, source.Source, primaryKey, objectType).Scan(db.SelectValues(rpslObject)...)
	if err != nil {
		return err
	}
	rpslObject.ToVersion = file.Version
	return db.Update(tx, rpslObject)
}

func selectObjectQuery() string {
	srcDesc := db.GetDescriptor(&pgpersist.NRTMSource{})
	rpslObjectDesc := db.GetDescriptor(&pgpersist.RPSLObject{})
	return fmt.Sprintf(`
		SELECT %v
		FROM %v
		JOIN %v ON %v.id = %v.nrtm_source_id
		WHERE
			%v.source ILIKE($1)
			AND UPPER(%v.primary_key) = UPPER($2)
			AND %v.object_type = UPPER($3)
			AND %v.to_version = 0`,
		strings.Join(rpslObjectDesc.ColumnNamesWithAlias(), ", "),
		rpslObjectDesc.TableNameWithAlias(),
		srcDesc.TableNameWithAlias(),
		srcDesc.TableAlias(),
		rpslObjectDesc.TableAlias(),
		srcDesc.TableAlias(),
		rpslObjectDesc.TableAlias(),
		rpslObjectDesc.TableAlias(),
		rpslObjectDesc.TableAlias(),
	)
}

func asNotification(n pgpersist.Notification) persist.Notification {
	return persist.Notification{
		ID:           n.ID,
		Version:      n.Version,
		NRTMSourceID: n.NRTMSourceID,
		Payload:      n.Payload,
		Created:      n.Created,
	}
}
