package persist

import (
	"time"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/pg/db"
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
