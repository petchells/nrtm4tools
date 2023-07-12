package persist

import (
	"time"

	"github.com/jackc/pgx/v5"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/nrtm4model"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/pg/db"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/rpsl"
)

type PgRepository struct {
}

func (repo PgRepository) InitializeConnectionPool(dbUrl string) {
	db.InitializeConnectionPool(dbUrl)
}

func (repo PgRepository) SaveState(state persist.NRTMState) *persist.ErrNrtmClient {
	var dbstate *NRTMState
	dbstate.Created = time.Now()

	err := db.WithTransaction(func(tx pgx.Tx) error {
		return nil
	})
	if err != nil {
		return &persist.ErrNrtmClient{Msg: "SaveState transaction failed"}
	}
	return nil
}

func (repo PgRepository) SaveSnapshotFile(state persist.NRTMState, snapshotObject nrtm4model.SnapshotFile) *persist.ErrNrtmClient {
	return nil
}

func (repo PgRepository) SaveSnapshotObject(state persist.NRTMState, snapshotObject nrtm4model.SnapshotObject) *persist.ErrNrtmClient {
	err := db.WithTransaction(func(tx pgx.Tx) error {
		var dbstate *NRTMState
		var err error
		err = db.GetByID(tx, state.ID, dbstate)
		if err != nil {
			return err
		}
		now := time.Now()
		rpslObject, err := rpsl.ParseString(snapshotObject.Object)
		if err != nil {
			return err
		}
		rpslObjectDB := RPSLObject{
			ID:          0,
			ObjectType:  rpslObject.ObjectType,
			RPSL:        rpslObject.Payload,
			Source:      state.Source,
			PrimaryKey:  rpslObject.PrimaryKey,
			NrtmstateID: dbstate.ID,
			Created:     now,
			Updated:     now,
		}
		return db.Create(tx, &rpslObjectDB)
	})
	if err != nil {
		clientErr, ok := err.(*persist.ErrNrtmClient)
		if ok {
			return clientErr
		}
		return &persist.ErrNrtmClient{Msg: err.Error()}
	}
	return nil
}

func (repo PgRepository) GetState(source string) (persist.NRTMState, *persist.ErrNrtmClient) {
	var state persist.NRTMState
	var dbstate *NRTMState
	err := db.WithTransaction(func(tx pgx.Tx) error {
		dbstate = GetLastState(tx, source)
		if dbstate == nil {
			return &persist.ErrNoState
		}
		return nil
	})
	if err != nil {
		clientErr, ok := err.(*persist.ErrNrtmClient)
		if ok {
			return state, clientErr
		}
		return state, &persist.ErrNrtmClient{Msg: err.Error()}
	}
	state.ID = dbstate.ID
	state.Created = dbstate.Created
	state.Source = dbstate.Source
	state.Version = dbstate.Version
	state.URL = dbstate.URL
	state.Type, _ = persist.ToFileType(dbstate.Type)
	state.FileName = dbstate.FileName
	return state, nil
}
