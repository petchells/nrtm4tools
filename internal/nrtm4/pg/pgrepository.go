package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/petchells/nrtm4tools/internal/nrtm4/persist"
	"github.com/petchells/nrtm4tools/internal/nrtm4/pg/db"
	pgpersist "github.com/petchells/nrtm4tools/internal/nrtm4/pg/persist"
	"github.com/petchells/nrtm4tools/internal/nrtm4/rpsl"
	"github.com/petchells/nrtm4tools/internal/nrtm4/util"
)

// PostgresRepository implementation of the Repository interface
type PostgresRepository struct {
}

// Initialize implementation of the Repository interface
func (repo PostgresRepository) Initialize(dbURL string) error {
	return db.InitializeConnectionPool(dbURL)
}

// ListSources returns a list of all sources
func (repo PostgresRepository) ListSources() ([]persist.NRTMSource, error) {
	var sources []persist.NRTMSource
	var err error
	var pgsources []pgpersist.NRTMSource
	err = db.WithTransaction(func(tx pgx.Tx) error {
		pgsources, err = db.GetAll(tx, pgpersist.NRTMSource{}, nil)
		return err
	})
	if err != nil {
		logger.Error("Error in ListSources", "error", err)
		return sources, err
	}
	for _, s := range pgsources {
		sources = append(sources, s.AsNRTMSource())
	}
	return sources, nil
}

// RemoveSource removes a source from the repo, including history
//
// A lock is put on table `nrtm_rpslobject` before deletions happen, so this might slow down
// updates to other sources when this is running.
func (repo PostgresRepository) RemoveSource(source persist.NRTMSource) error {
	err := db.WithTransaction(func(tx pgx.Tx) error {
		type pgcmd struct {
			sql  string
			args []any
		}
		cmds := []pgcmd{
			{`
			DELETE FROM
				nrtm_rpslobject_history
			WHERE source_id = $1
			`, []any{source.ID},
			}, {`
			DELETE FROM
				nrtm_notification
			WHERE source_id = $1
			`, []any{source.ID},
			}, {`
			LOCK TABLE nrtm_rpslobject IN SHARE MODE
			`, []any{},
			}, {`
			ALTER TABLE
				nrtm_rpslobject
			DISABLE TRIGGER modify_rpsl_trigger
			`, []any{},
			}, {`
			DELETE FROM
				nrtm_rpslobject
			WHERE source_id = $1
			`, []any{source.ID},
			}, {`
			ALTER TABLE
				nrtm_rpslobject
			ENABLE TRIGGER modify_rpsl_trigger
			`, []any{},
			}, {`
			DELETE FROM
				nrtm_source
			WHERE id = $1
			`, []any{source.ID},
			},
		}
		for _, cmd := range cmds {
			_, err := tx.Exec(context.Background(), cmd.sql, cmd.args...)
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
		WHERE source_id = $1
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
			err = rows.Scan(db.ValuesForSelect(&ent)...)
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
				db.NextID(),
				rpslObject.ObjectType,
				rpslObject.PrimaryKey,
				source.ID,
				file.Version,
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
			for _, obj := range rpslObjects {
				types.Add(obj.ObjectType)
				logger.Debug("Possible failure", "type", obj.ObjectType, "primaryKey", obj.PrimaryKey)
			}
			logger.Warn("Failed to save objects", "types", types.String(), "error", err)
			return err
		}
		return nil
	})
}

// AddModifyObject updates an RPSL finding the current matching pk then updating or adding
func (repo PostgresRepository) AddModifyObject(
	source persist.NRTMSource,
	rpsl rpsl.Rpsl,
	file persist.NrtmFileJSON,
) error {
	newRow := &pgpersist.RPSLObject{
		ObjectType: rpsl.ObjectType,
		PrimaryKey: rpsl.PrimaryKey,
		SourceID:   source.ID,
		Version:    uint32(file.Version),
		RPSL:       rpsl.Payload,
	}
	return db.WithTransaction(func(tx pgx.Tx) error {

		var err error

		sql := selectCurrentObjectQuery()
		rpslObject := new(pgpersist.RPSLObject)
		err = tx.QueryRow(context.Background(), sql, source.ID, rpsl.PrimaryKey, rpsl.ObjectType).Scan(db.ValuesForSelect(rpslObject)...)
		if err != nil && err != pgx.ErrNoRows {
			return err
		}
		if err == pgx.ErrNoRows {
			newRow.ID = db.NextID()
			return db.Create(tx, newRow)
		}
		newRow.ID = rpslObject.ID
		return db.Update(tx, newRow)
	})
}

// DeleteObject removes a row matching the params
func (repo PostgresRepository) DeleteObject(
	source persist.NRTMSource,
	objectType string,
	primaryKey string,
	file persist.NrtmFileJSON,
) error {
	return db.WithTransaction(func(tx pgx.Tx) error {
		sql := selectCurrentObjectQuery()
		rpslObject := new(pgpersist.RPSLObject)
		err := tx.QueryRow(context.Background(), sql, source.ID, primaryKey, objectType).Scan(db.ValuesForSelect(rpslObject)...)
		if err != nil {
			return err
		}
		rpslObjectDesc := db.GetDescriptor(&pgpersist.RPSLObject{})
		sql = fmt.Sprintf(`DELETE FROM %v WHERE id=$1`, rpslObjectDesc.TableName())
		_, err = tx.Exec(context.Background(), sql, rpslObject.ID)
		return err
	})
}

func selectCurrentObjectQuery() string {
	rpslObjectDesc := db.GetDescriptor(&pgpersist.RPSLObject{})
	return fmt.Sprintf(`
		SELECT %v
		FROM %v
		WHERE
			source_id = $1
			AND primary_key = UPPER($2)
			AND object_type = UPPER($3)`,
		rpslObjectDesc.ColumnNamesCommaSeparated(),
		rpslObjectDesc.TableName(),
	)
}

func asNotification(n pgpersist.Notification) persist.Notification {
	return persist.Notification{
		ID:       n.ID,
		Version:  n.Version,
		SourceID: n.SourceID,
		Payload:  n.Payload,
		Created:  n.Created,
	}
}
