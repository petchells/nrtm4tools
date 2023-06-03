package persist

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/pg/db"
)

type PgRepository struct {
}

func (repo PgRepository) InitializeConnectionPool(dbUrl string) {
	db.InitializeConnectionPool(dbUrl)
}

func (repo PgRepository) GetState(source string) (persist.NRTMState, error) {
	var state persist.NRTMState
	var dbstate *NRTMState
	err := db.WithTransaction(func(tx pgx.Tx) error {
		dbstate = GetLastState(tx, source)
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
	state.Source = dbstate.Source
	state.Version = dbstate.Version
	state.IsDelta = dbstate.IsDelta
	state.URL = dbstate.URL
	state.Delta = dbstate.Delta
	state.SnapshotPath = dbstate.SnapshotPath
	return state, nil
}
