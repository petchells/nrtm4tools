package pg

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/nrtm4model"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/pg/db"
	pgpersist "gitlab.com/etchells/nrtm4client/internal/nrtm4/pg/persist"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/rpsl"
)

// PostgresRepository implementation of the Repository interface
type PostgresRepository struct {
	persist.Repository
}

// Initialize implementation of the Repository interface
func (repo *PostgresRepository) Initialize(dbURL string) error {
	return db.InitializeConnectionPool(dbURL)
}

// Close implementation of the interface. Nothing needed for pg (for now)
func (repo *PostgresRepository) Close() error {
	return nil
}

// SaveState saves a reference to an NRTM file
func (repo *PostgresRepository) SaveState(state *persist.NRTMState) error {
	return db.WithTransaction(func(tx pgx.Tx) error {
		st := pgpersist.NRTMState{
			ID:       uint64(db.NextID()),
			Source:   state.Source,
			Version:  state.Version,
			URL:      state.URL,
			Type:     persist.NotificationFile.String(),
			FileName: state.FileName,
			Created:  time.Now().UTC(),
		}
		state.ID = st.ID
		return db.Create(tx, &st)
	})
}

// SaveSnapshotFile not supported
func (repo *PostgresRepository) SaveSnapshotFile(state persist.NRTMState, snapshotObject nrtm4model.SnapshotFile) error {
	return nil
}

// SaveSnapshotObject save am RPSL object
func (repo *PostgresRepository) SaveSnapshotObject(state persist.NRTMState, rpslObject rpsl.Rpsl) error {
	return db.WithTransaction(func(tx pgx.Tx) error {
		dbstate := new(pgpersist.NRTMState)
		err := db.GetByID(tx, state.ID, dbstate)
		if err != nil {
			return err
		}
		now := time.Now().UTC()
		rpslObjectDB := pgpersist.RPSLObject{
			ID:          uint64(db.NextID()),
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
}

// SaveSnapshotObjects saves a list of rpsl object in a go routine
func (repo *PostgresRepository) SaveSnapshotObjects(state persist.NRTMState, rpslObjects []rpsl.Rpsl) error {

	ch := make(chan error)
	updateDB := func(ch chan error) {
		err := db.WithTransaction(func(tx pgx.Tx) error {
			dbstate := new(pgpersist.NRTMState)
			err := db.GetByID(tx, state.ID, dbstate)
			if err != nil {
				return err
			}
			now := time.Now().UTC()
			inputRows := [][]any{}
			for _, rpslObject := range rpslObjects {
				inputRow := []any{
					uint64(db.NextID()),
					rpslObject.Source,
					rpslObject.ObjectType,
					rpslObject.PrimaryKey,
					rpslObject.Payload,
					dbstate.ID,
					now,
					now,
				}
				inputRows = append(inputRows, inputRow)
			}
			rpslDescriptor := db.GetDescriptor(&pgpersist.RPSLObject{})
			_, err = tx.CopyFrom(context.Background(), pgx.Identifier{rpslDescriptor.TableName()}, rpslDescriptor.ColumnNames(), pgx.CopyFromRows(inputRows))
			if err != nil {
				log.Println("WARNING failed to save objects with error", err)
				return err
			}
			return nil
		})
		ch <- err
	}
	go updateDB(ch)
	return <-ch
}

// GetState get the last known state for the source
func (repo *PostgresRepository) GetState(source string) (persist.NRTMState, error) {
	var state persist.NRTMState
	var dbstate *pgpersist.NRTMState
	err := db.WithTransaction(func(tx pgx.Tx) error {
		dbstate = pgpersist.GetLastState(tx, source)
		if dbstate == nil {
			return &persist.ErrStateNotInitialized
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
	state.FileName = dbstate.FileName
	return state, nil
}
