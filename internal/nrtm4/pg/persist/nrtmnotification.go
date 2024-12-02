package persist

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/pg/db"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/util"
)

// Notification is a binding to a PG database table
type Notification struct {
	db.EntityManaged `em:"nrtm_notification nnot"`
	ID               uint64                   `em:"."`
	Version          uint32                   `em:"."`
	NRTMSourceID     uint64                   `em:"."`
	Payload          persist.NotificationJSON `em:"."`
	Created          time.Time                `em:"."`
}

// NewNotification saves a notification in the database
func NewNotification(tx pgx.Tx, sourceID uint64, payload persist.NotificationJSON) error {
	lastN := new(Notification)
	descr := db.GetDescriptor(lastN)
	sql := fmt.Sprintf(`
		SELECT %v
		FROM %v
		WHERE nrtm_source_id = $1
		ORDER BY
			version DESC,
			created DESC
		LIMIT 1
	`, descr.ColumnNamesCommaSeparated(), descr.TableName())

	newNotification := func(tx pgx.Tx) error {
		logger.Debug("Saving new notification")
		return db.Create(tx, &Notification{
			ID:           db.NextID(),
			Version:      payload.Version,
			NRTMSourceID: sourceID,
			Payload:      payload,
			Created:      util.AppClock.Now(),
		})
	}

	err := tx.QueryRow(context.Background(), sql, sourceID).Scan(db.SelectValues(lastN)...)
	if err == pgx.ErrNoRows {
		return newNotification(tx)
	} else if err != nil {
		return err
	}
	if payload.Version == lastN.Version {
		if payload.SnapshotRef.Version == lastN.Payload.SnapshotRef.Version {
			// Nothing to do
			return nil
		}
		return newNotification(tx)
	} else if payload.Version > lastN.Version {
		return newNotification(tx)
	}
	logger.Error("Expected higher notification version", "lastN.Version", lastN.Version, "payload.Version", payload.Version)
	return errors.New("expected higher notification version than the one found")
}
