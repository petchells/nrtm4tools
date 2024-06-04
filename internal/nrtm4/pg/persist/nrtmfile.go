package persist

import (
	"time"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/pg/db"
)

// NRTMFile is a binding to a PG database table
type NRTMFile struct {
	db.EntityManaged `em:"nrtm_file nf"`
	ID               uint64    `em:"."`
	Created          time.Time `em:"."`
	FileName         string    `em:"."`
	NRTMSourceID     uint64    `em:"."`
	Type             string    `em:"."`
	URL              string    `em:"."`
	Version          uint32    `em:"."`
}
