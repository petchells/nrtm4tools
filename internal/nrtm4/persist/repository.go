package persist

import (
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/rpsl"
)

// Repository defines the functions for NRTMClient's persistent storage
type Repository interface {
	Initialize(dbURL string) error
	SaveSource(NRTMSource) (NRTMSource, error)
	GetSources() ([]NRTMSource, error)
	SaveFile(*NRTMFile) error
	SaveSnapshotObject(NRTMSource, rpsl.Rpsl) error
	SaveSnapshotObjects(NRTMSource, []rpsl.Rpsl) error
	Close() error
}
