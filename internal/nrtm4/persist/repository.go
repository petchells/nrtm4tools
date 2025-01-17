package persist

import (
	"github.com/petchells/nrtm4client/internal/nrtm4/rpsl"
)

// Repository defines the functions for NRTMClient's persistent storage
type Repository interface {
	Initialize(string) error
	SaveSource(NRTMSource, NotificationJSON) (NRTMSource, error)
	RemoveSource(NRTMSource) error
	GetSources() ([]NRTMSource, error)
	GetNotificationHistory(NRTMSource, uint32, uint32) ([]Notification, error)
	SaveFile(*NRTMFile) error
	SaveSnapshotObjects(NRTMSource, []rpsl.Rpsl, NrtmFileJSON) error
	AddModifyObject(NRTMSource, rpsl.Rpsl, NrtmFileJSON) error
	DeleteObject(NRTMSource, string, string, NrtmFileJSON) error
	Close() error
}
