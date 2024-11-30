package persist

import (
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
func NewNotification(tx pgx.Tx, version uint32, sourceID uint64, payload persist.NotificationJSON) (*Notification, error) {
	n := Notification{
		ID:           db.NextID(),
		Version:      version,
		NRTMSourceID: sourceID,
		Payload:      payload,
		Created:      util.AppClock.Now(),
	}
	if err := db.Create(tx, &n); err != nil {
		return nil, err
	}
	return &n, nil
}
