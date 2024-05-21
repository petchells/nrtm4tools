package persist

import (
	"errors"
	"strings"
	"time"
)

// NRTMState describes a downloaded NRTM file
type NRTMState struct {
	ID       uint64
	Created  time.Time
	Source   string
	Version  uint
	URL      string
	Type     NTRMFileType
	FileName string
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

func ToFileType(s string) (NTRMFileType, error) {
	target := strings.ToLower(s)
	for i, str := range ftstrings {
		if str == target {
			return NTRMFileType(i), nil
		}
	}
	return -1, errors.New("out of range")
}
