package persist

import (
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/pg/db"
)

// RPSLObject is an RPSL object
type RPSLObject struct {
	db.EntityManaged `em:"nrtm_rpslobject ro"`
	ID               uint64 `em:"."`
	ObjectType       string `em:"."`
	PrimaryKey       string `em:"."`
	RPSL             string `em:"."`
	NRTMSourceID     uint64 `em:"."`
	FromVersion      uint32 `em:"."`
	ToVersion        uint32 `em:"."`
}
