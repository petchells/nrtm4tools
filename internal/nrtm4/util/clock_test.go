package util

import (
	"log"
	"testing"
	"time"
)

var (
	dateStr = "2020-02-02T02:02:02.987654321Z"
)

func NewTestClock(startTime string) Clock {
	testTimestamp, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		log.Fatalln("Failed to parse timestamp", startTime)
	}
	timeOffset := time.Now().UTC().UnixNano() - testTimestamp.UnixNano()
	return testClock{time.Duration(timeOffset * int64(time.Nanosecond))}
}

type testClock struct {
	offset time.Duration
}

func (tc testClock) Now() time.Time {
	return time.Now().UTC().Add(-1 * tc.offset)
}

func (tc testClock) After(d time.Duration) <-chan time.Time {
	time.Sleep(d)
	c := make(chan time.Time, 1)
	c <- tc.Now()
	return c
}

func TestTimeFormat(t *testing.T) {
	AppClock = NewTestClock(dateStr)
	ts := AppClock.Now().Format(RFC3339Milli)
	expected := "2020-02-02T02:02:02.987Z"
	if ts != expected {
		t.Error("Format failed expected", expected, "but was", ts)
	}
}
