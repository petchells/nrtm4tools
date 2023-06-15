package persist

import (
	"time"

	"github.com/jackc/pgx/v5"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/pg/db"
)

type PgRepository struct {
}

func (repo PgRepository) InitializeConnectionPool(dbUrl string) {
	db.InitializeConnectionPool(dbUrl)
}

func (repo PgRepository) SaveState(state persist.NRTMState) error {
	var dbstate *NRTMState
	dbstate.Created = time.Now()

	err := db.WithTransaction(func(tx pgx.Tx) error {
		return nil
	})
	return err
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
	state.URL = dbstate.URL
	state.Type, _ = persist.ToFileType(dbstate.Type)
	state.Payload = dbstate.Payload
	return state, nil
}
