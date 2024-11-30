package persist

import (
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/rpsl"
)

// Repository defines the functions for NRTMClient's persistent storage
type Repository interface {
	Initialize(dbURL string) error
	SaveSource(NRTMSource, NotificationJSON) (NRTMSource, error)
	GetSources() ([]NRTMSource, error)
	SaveFile(*NRTMFile) error
	SaveNotification(NRTMSource, Notification) error
	SaveSnapshotObjects(NRTMSource, []rpsl.Rpsl, NrtmFileJSON) error
	AddModifyObject(NRTMSource, rpsl.Rpsl, NrtmFileJSON) error
	DeleteObject(NRTMSource, string, string, NrtmFileJSON) error
	Close() error
}
