package persist

import (
	"errors"
	"strings"
	"time"
)

// NRTMSource holds information about a remote NRTM source
type NRTMSource struct {
	ID              uint64 `json:",string"`
	Source          string
	SessionID       string
	Version         uint32
	NotificationURL string
	Label           string
	Created         time.Time
}

// NRTMSourceDetails is a source with notification objects
type NRTMSourceDetails struct {
	NRTMSource
	Notifications []Notification
}

// NewNRTMSource prepares a new source object
func NewNRTMSource(notification NotificationJSON, label string, notificationURL string) NRTMSource {
	return NRTMSource{
		Source:          notification.Source,
		SessionID:       notification.SessionID,
		Version:         notification.SnapshotRef.Version,
		Label:           label,
		NotificationURL: notificationURL,
	}
}

// Notification is a domain representation of a notification file
type Notification struct {
	ID           uint64 `json:",string"`
	Version      uint32
	NRTMSourceID uint64 `json:",string"`
	Payload      NotificationJSON
	Created      time.Time
}

// NRTMFile describes a downloaded NRTM file
type NRTMFile struct {
	ID           uint64 `json:",string"`
	Version      uint32
	Type         NTRMFileType
	URL          string
	FileName     string
	NrtmSourceID uint64 `json:",string"`
	Created      time.Time
}

// NTRMFileType enumerator for file types
type NTRMFileType int

const (

	// NotificationFile notification file
	NotificationFile NTRMFileType = iota

	// SnapshotFile snapshot file
	SnapshotFile

	// DeltaFile delta file
	DeltaFile
)

var ftstrings = [...]string{"notification", "snapshot", "delta"}

func (ft NTRMFileType) String() string {
	if ft < NotificationFile || ft > DeltaFile {
		return ""
	}
	return ftstrings[ft]
}

// ToFileType returns an NRTMFileType which matches s
func ToFileType(s string) (NTRMFileType, error) {
	target := strings.ToLower(s)
	for i := range ftstrings {
		if target == ftstrings[i] {
			return NTRMFileType(i), nil
		}
	}
	return -1, errors.New("out of range")
}
