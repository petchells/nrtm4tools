package persist

import (
	"time"

	"github.com/petchells/nrtm4tools/internal/nrtm4/persist"
	"github.com/petchells/nrtm4tools/internal/nrtm4/pg/db"
	"github.com/petchells/nrtm4tools/internal/nrtm4/util"
)

// NRTMSource pg database mapping for nrtm_source
type NRTMSource struct {
	db.EntityManaged `em:"nrtm_source src"`
	ID               uint64    `em:"-"`
	Source           string    `em:"-"`
	SessionID        string    `em:"-"`
	Version          uint32    `em:"-"`
	NotificationURL  string    `em:"-"`
	Label            string    `em:"-"`
	Created          time.Time `em:"-"`
}

// NewNRTMSource is a shorthand function which prepares a source object for storage
func NewNRTMSource(source persist.NRTMSource) NRTMSource {
	id := db.NextID()
	sourceObj := NRTMSource{
		ID:              id,
		Source:          source.Source,
		SessionID:       source.SessionID,
		Version:         source.Version,
		NotificationURL: source.NotificationURL,
		Label:           source.Label,
		Created:         util.AppClock.Now(),
	}
	return sourceObj
}

// FromNRTMSource is a shorthand function which transforms a Pg source to a generic persist source
func FromNRTMSource(source persist.NRTMSource) NRTMSource {
	return NRTMSource{
		ID:              source.ID,
		Source:          source.Source,
		SessionID:       source.SessionID,
		Version:         source.Version,
		NotificationURL: source.NotificationURL,
		Label:           source.Label,
		Created:         source.Created,
	}
}

// AsNRTMSource return this row as a app-level source
func (s *NRTMSource) AsNRTMSource() persist.NRTMSource {
	return persist.NRTMSource{
		ID:              s.ID,
		Source:          s.Source,
		SessionID:       s.SessionID,
		Version:         s.Version,
		NotificationURL: s.NotificationURL,
		Label:           s.Label,
		Created:         s.Created,
	}
}
