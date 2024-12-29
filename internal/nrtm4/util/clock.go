package util

import "time"

// AppClock is the application's stubbable clock, set to UTC
var AppClock Clock = realClock{}

// RFC3339Milli format that rounds to ms, so browsers can grok it
const RFC3339Milli = "2006-01-02T15:04:05.999Z07:00"

// Clock defies functions of a clocky nature
type Clock interface {
	Now() time.Time
	After(d time.Duration) <-chan time.Time
}

type realClock struct{}

func (realClock) Now() time.Time                         { return time.Now().UTC() }
func (realClock) After(d time.Duration) <-chan time.Time { return time.After(d) }
