package util

import "time"

// AppClock is the application's stubbable clock, set to UTC
var AppClock Clock = realClock{}

// Clock defines what we want from an application clock
type Clock interface {
	Now() time.Time
	After(d time.Duration) <-chan time.Time
}

type realClock struct{}

func (realClock) Now() time.Time                         { return time.Now().UTC() }
func (realClock) After(d time.Duration) <-chan time.Time { return time.After(d) }
