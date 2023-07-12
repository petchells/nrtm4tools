package persist

import (
	"time"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/pg/db"
)

type RPSLObject struct {
	db.EntityManaged `em:"rpslobject ro"`
	ID               uint64    `em:"."`
	ObjectType       string    `em:"."`
	RPSL             string    `em:"."`
	Source           string    `em:"."`
	PrimaryKey       string    `em:"."`
	NrtmstateID      uint64    `em:"."`
	Created          time.Time `em:"."`
	Updated          time.Time `em:"."`
}
