package persist

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/pg/db"
)

type PgRepository struct {
	persist.Repository
}

func (repo PgRepository) InitializeConnectionPool(dbUrl string) {
	db.InitializeConnectionPool(dbUrl)
}

func (repo PgRepository) GetState() (persist.NRTMState, error) {
	var state persist.NRTMState
	var dbstate *NRTMState
	err := db.WithTransaction(func(tx pgx.Tx) error {
		dbstate = GetLastState(tx)
		if dbstate == nil {
			return persist.ErrNoState
		}
		return nil
	})
	if err != nil {
		return state, err
	}
	state.ID = dbstate.ID
	state.Created = dbstate.Created
	state.URL = dbstate.URL
	state.Version = dbstate.FileVersion
	state.IsDelta = dbstate.IsDelta
	state.Delta = dbstate.Delta
	state.SnapshotPath = dbstate.SnapshotPath
	return state, nil
}
