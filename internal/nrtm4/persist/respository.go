package persist

import (
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/nrtm4model"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/rpsl"
)

type Repository interface {
	Initialize(dbUrl string) error
	GetState(string) (NRTMState, error)
	SaveState(*NRTMState) error
	SaveSnapshotFile(NRTMState, nrtm4model.SnapshotFile) error
	SaveSnapshotObject(NRTMState, rpsl.Rpsl) error
	SaveSnapshotObjects(NRTMState, []rpsl.Rpsl) error
	Close() error
}
