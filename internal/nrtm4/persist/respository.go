package persist

import (
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/nrtm4model"
)

type Repository interface {
	InitializeConnectionPool(dbUrl string)
	GetState(string) (NRTMState, *ErrNrtmClient)
	SaveState(NRTMState) *ErrNrtmClient
	SaveSnapshotFile(NRTMState, nrtm4model.SnapshotFile) *ErrNrtmClient
	SaveSnapshotObject(NRTMState, nrtm4model.SnapshotObject) *ErrNrtmClient
}
