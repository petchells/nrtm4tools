package persist

import "time"

type NRTMState struct {
	ID           uint64
	Created      time.Time
	URL          string
	Version      uint
	IsDelta      bool
	Delta        string
	SnapshotPath string
}
