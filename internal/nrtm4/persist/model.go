package persist

import (
	"errors"
	"strings"
	"time"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/nrtm4model"
)

// NRTMSource holds information about a remote NRTM source
type NRTMSource struct {
	ID              uint64
	Source          string
	SessionID       string
	Version         uint32
	NotificationURL string
	Label           string
	Created         time.Time
}

// NewNRTMSource prepares a new source object
func NewNRTMSource(notification nrtm4model.NotificationJSON, label string, notificationURL string) NRTMSource {
	return NRTMSource{
		Source:          notification.Source,
		SessionID:       notification.SessionID,
		Version:         notification.Version,
		Label:           label,
		NotificationURL: notificationURL,
	}
}

// NRTMFile describes a downloaded NRTM file
type NRTMFile struct {
	ID           uint64
	Version      uint
	Type         NTRMFileType
	URL          string
	FileName     string
	NrtmSourceID uint64
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
	for i, str := range ftstrings {
		if str == target {
			return NTRMFileType(i), nil
		}
	}
	return -1, errors.New("out of range")
}
