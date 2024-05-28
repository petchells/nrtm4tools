package pg

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/pg/db"
	pgpersist "gitlab.com/etchells/nrtm4client/internal/nrtm4/pg/persist"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/rpsl"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/util"
)

// PostgresRepository implementation of the Repository interface
type PostgresRepository struct {
}

// Initialize implementation of the Repository interface
func (repo *PostgresRepository) Initialize(dbURL string) error {
	return db.InitializeConnectionPool(dbURL)
}

// GetSources returns a list of all sources
func (repo *PostgresRepository) GetSources() ([]persist.NRTMSource, error) {
	var sources []persist.NRTMSource
	var err error
	var pgsources []pgpersist.NRTMSource
	src := new(pgpersist.NRTMSource)
	err = db.WithTransaction(func(tx pgx.Tx) error {
		pgsources, err = db.GetAll(tx, *src, nil)
		return err
	})
	if err != nil {
		logger.Error("Error in GetSources", err)
		return sources, err
	}
	for _, s := range pgsources {
		sources = append(sources, s.AsNRTMSource())
	}
	return sources, nil
}

// SaveSource updates a source if ID is non-zero, or creates a new one if it is
func (repo *PostgresRepository) SaveSource(source persist.NRTMSource) (persist.NRTMSource, error) {
	var pgSource pgpersist.NRTMSource
	err := db.WithTransaction(func(tx pgx.Tx) error {
		if source.ID == 0 {
			pgSource = pgpersist.NewNRTMSource(source)
			return db.Create(tx, &pgSource)
		}
		pgSource = pgpersist.FromNRTMSource(source)
		return db.Update(tx, &pgSource)
	})
	return pgSource.AsNRTMSource(), err
}

// Close implementation of the interface. Nothing needed for pg (for now)
func (repo *PostgresRepository) Close() error {
	return nil
}

// SaveFile saves a reference to an NRTM file
func (repo *PostgresRepository) SaveFile(nrtmFile *persist.NRTMFile) error {
	return db.WithTransaction(func(tx pgx.Tx) error {
		st := pgpersist.NRTMFile{
			ID:           uint64(db.NextID()),
			Version:      nrtmFile.Version,
			URL:          nrtmFile.URL,
			Type:         nrtmFile.Type.String(),
			NRTMSourceID: nrtmFile.NrtmSourceID,
			FileName:     nrtmFile.FileName,
			Created:      time.Now().UTC(),
		}
		nrtmFile.ID = st.ID
		return db.Create(tx, &st)
	})
}

// SaveSnapshotObject save am RPSL object
func (repo *PostgresRepository) SaveSnapshotObject(source persist.NRTMSource, rpslObject rpsl.Rpsl) error {
	return db.WithTransaction(func(tx pgx.Tx) error {
		rpslObjectDB := pgpersist.RPSLObject{
			ID:           uint64(db.NextID()),
			ObjectType:   rpslObject.ObjectType,
			RPSL:         rpslObject.Payload,
			PrimaryKey:   rpslObject.PrimaryKey,
			NRTMSourceID: source.ID,
			FromVersion:  source.Version,
		}
		return db.Create(tx, &rpslObjectDB)
	})
}

// SaveSnapshotObjects saves a list of rpsl object in a go routine
func (repo *PostgresRepository) SaveSnapshotObjects(source persist.NRTMSource, rpslObjects []rpsl.Rpsl) error {

	ch := make(chan error)
	updateDB := func(ch chan error) {
		err := db.WithTransaction(func(tx pgx.Tx) error {
			var err error
			inputRows := [][]any{}
			for _, rpslObject := range rpslObjects {
				inputRow := []any{
					uint64(db.NextID()),
					rpslObject.ObjectType,
					rpslObject.PrimaryKey,
					rpslObject.Payload,
					source.ID,
					source.Version,
					0,
				}
				inputRows = append(inputRows, inputRow)
			}
			rpslDescriptor := db.GetDescriptor(&pgpersist.RPSLObject{})
			_, err = tx.CopyFrom(context.Background(), pgx.Identifier{rpslDescriptor.TableName()}, rpslDescriptor.ColumnNames(), pgx.CopyFromRows(inputRows))
			if err != nil {
				types := util.NewSet[string]()
				for _, inp := range rpslObjects {
					types.Add(inp.ObjectType)
				}
				logger.Warn("Failed to save objects", "types", types.String(), err)
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
func (repo *PostgresRepository) GetState(source string) (persist.NRTMFile, error) {
	var state persist.NRTMFile
	var dbstate *pgpersist.NRTMFile
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
	state.Version = dbstate.Version
	state.URL = dbstate.URL
	state.Type, _ = persist.ToFileType(dbstate.Type)
	state.FileName = dbstate.FileName
	return state, nil
}
