package util

import (
	"testing"
	"time"
)

var (
	dateStr    = "2020-02-02T02:02:02.987654321Z"
	fakeNow, _ = time.Parse(time.RFC3339, dateStr)
)

type testClock struct {
	fakeNow time.Time
}

func (tc testClock) Now() time.Time {
	return tc.fakeNow
}

func (tc testClock) After(d time.Duration) <-chan time.Time {
	time.Sleep(d)
	delayed := tc.fakeNow.Add(d)
	c := make(chan time.Time, 1)
	c <- delayed
	return c
}

func TestTimeFormat(t *testing.T) {
	AppClock = testClock{fakeNow: fakeNow}
	ts := AppClock.Now().Format(RFC3339Milli)
	expected := "2020-02-02T02:02:02.987Z"
	if ts != expected {
		t.Error("Format failed expected", expected, "but was", ts)
	}
}
