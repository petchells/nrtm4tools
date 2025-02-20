package persist

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/petchells/nrtm4client/internal/nrtm4/persist"
	"github.com/petchells/nrtm4client/internal/nrtm4/pg/db"
	"github.com/petchells/nrtm4client/internal/nrtm4/util"
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

	pver := uint32(payload.Version)
	newNotification := func(tx pgx.Tx) error {
		logger.Debug("Saving new notification")
		return db.Create(tx, &Notification{
			ID:           db.NextID(),
			Version:      pver,
			NRTMSourceID: sourceID,
			Payload:      payload,
			Created:      util.AppClock.Now(),
		})
	}

	err := tx.QueryRow(context.Background(), sql, sourceID).Scan(db.ValuesForSelect(lastN)...)
	if err == pgx.ErrNoRows {
		return newNotification(tx)
	} else if err != nil {
		return err
	}
	if pver == lastN.Version {
		if payload.SnapshotRef.Version == lastN.Payload.SnapshotRef.Version {
			// Nothing to do
			return nil
		}
		return newNotification(tx)
	} else if pver > lastN.Version {
		return newNotification(tx)
	}
	logger.Error("Expected higher notification version", "lastN.Version", lastN.Version, "payload.Version", pver)
	return errors.New("expected higher notification version than the one found")
}
