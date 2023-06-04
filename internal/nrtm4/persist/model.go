package persist

import "time"

type NRTMState struct {
	ID      uint64
	Created time.Time
	Source  string
	Version int
	URL     string
	Type    string
	Payload string
}
