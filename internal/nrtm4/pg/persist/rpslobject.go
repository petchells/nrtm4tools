package persist

import (
	"github.com/petchells/nrtm4tools/internal/nrtm4/pg/db"
)

// RPSLObject is an RPSL object
type RPSLObject struct {
	db.EntityManaged `em:"nrtm_rpslobject rpsl"`
	ID               uint64 `em:"-"`
	ObjectType       string `em:"-"`
	PrimaryKey       string `em:"-"`
	SourceID         uint64 `em:"-"`
	Version          uint32 `em:"-"`
	RPSL             string `em:"-"`
}
