package persist

import (
	"errors"
	"strings"
	"time"
)

type NRTMState struct {
	ID       uint64
	Created  time.Time
	Source   string
	Version  uint
	URL      string
	Type     NTRMFileType
	FileName string
}

type NTRMFileType int

const (
	NotificationFile NTRMFileType = iota
	SnapshotFile
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
