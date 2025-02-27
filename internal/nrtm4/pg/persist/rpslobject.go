package persist

import (
	"github.com/petchells/nrtm4tools/internal/nrtm4/pg/db"
)

// RPSLObject is an RPSL object
type RPSLObject struct {
	db.EntityManaged `em:"nrtm_rpslobject rpsl"`
	ID               uint64 `em:"."`
	ObjectType       string `em:"."`
	PrimaryKey       string `em:"."`
	NRTMSourceID     uint64 `em:"."`
	FromVersion      uint32 `em:"."`
	ToVersion        uint32 `em:"."`
	RPSL             string `em:"."`
}
