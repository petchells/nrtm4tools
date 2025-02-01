package pg

import (
	"context"
	"fmt"

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

// RemoveSource removes a source from the repo
func (repo PostgresRepository) RemoveSource(source persist.NRTMSource) error {
	err := db.WithTransaction(func(tx pgx.Tx) error {
		sqls := []string{`
			DELETE FROM
				nrtm_rpslobject
			WHERE nrtm_source_id = $1
			`, `
			DELETE FROM
				nrtm_notification
			WHERE nrtm_source_id = $1
			`, `
			DELETE FROM
				nrtm_file
			WHERE nrtm_source_id = $1
			`, `
			DELETE FROM
				nrtm_source
			WHERE id = $1
			`}
		for _, sql := range sqls {
			_, err := tx.Exec(context.Background(), sql, source.ID)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		logger.Error("Error in RemoveSource", "error", err)
		return err
	}
	return nil
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
		LIMIT 100
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
	newRow := &pgpersist.RPSLObject{
		ObjectType:   rpsl.ObjectType,
		PrimaryKey:   rpsl.PrimaryKey,
		NRTMSourceID: source.ID,
		FromVersion:  file.Version,
		RPSL:         rpsl.Payload,
	}
	return db.WithTransaction(func(tx pgx.Tx) error {

		var err error

		curDelta := getPossibleCurrentDeltaFrom(tx, *newRow)
		if curDelta != nil {
			// Already processed an operation, just overwrite it
			newRow.ID = curDelta.ID
			return db.Update(tx, newRow)
		}

		sql := selectCurrentObjectQuery()
		rpslObject := new(pgpersist.RPSLObject)
		err = tx.QueryRow(context.Background(), sql, source.ID, rpsl.PrimaryKey, rpsl.ObjectType).Scan(db.SelectValues(rpslObject)...)
		if err != nil && err != pgx.ErrNoRows {
			return err
		}
		if err != pgx.ErrNoRows {
			rpslObject.ToVersion = file.Version
			err = db.Update(tx, rpslObject)
			if err != nil {
				return err
			}
		}
		newRow.ID = db.NextID()
		return db.Create(tx, newRow)
	})
}

// DeleteObject doesn't remove any rows, instead it sets `to_version` to the file version
func (repo PostgresRepository) DeleteObject(
	source persist.NRTMSource,
	objectType string,
	primaryKey string,
	file persist.NrtmFileJSON,
) error {
	return db.WithTransaction(func(tx pgx.Tx) error {
		sql := selectCurrentObjectQuery()
		rpslObject := new(pgpersist.RPSLObject)
		err := tx.QueryRow(context.Background(), sql, source.ID, primaryKey, objectType).Scan(db.SelectValues(rpslObject)...)
		if err != nil {
			return err
		}
		rpslObject.ToVersion = file.Version
		return db.Update(tx, rpslObject)
	})
}

func selectCurrentObjectQuery() string {
	rpslObjectDesc := db.GetDescriptor(&pgpersist.RPSLObject{})
	return fmt.Sprintf(`
		SELECT %v
		FROM %v
		WHERE
			nrtm_source_id = $1
			AND primary_key = UPPER($2)
			AND object_type = UPPER($3)
			AND to_version = 0`,
		rpslObjectDesc.ColumnNamesCommaSeparated(),
		rpslObjectDesc.TableName(),
	)
}

func getPossibleCurrentDeltaFrom(tx pgx.Tx, curRPSL pgpersist.RPSLObject) *pgpersist.RPSLObject {
	rpslObjectDesc := db.GetDescriptor(&pgpersist.RPSLObject{})
	sql := fmt.Sprintf(`
		SELECT %v
		FROM %v
		WHERE
			nrtm_source_id = $1
			AND primary_key = UPPER($2)
			AND object_type = UPPER($3)
			AND from_version = $4`,
		rpslObjectDesc.ColumnNamesCommaSeparated(),
		rpslObjectDesc.TableName(),
	)
	rpslObject := new(pgpersist.RPSLObject)
	err := tx.QueryRow(context.Background(), sql, curRPSL.NRTMSourceID, curRPSL.PrimaryKey, curRPSL.ObjectType, curRPSL.FromVersion).Scan(db.SelectValues(rpslObject)...)
	if err != nil {
		if err != pgx.ErrNoRows {
			logger.Error("Could not get current delta", "curRPSL", curRPSL)
		}
		return nil
	}
	return rpslObject
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
