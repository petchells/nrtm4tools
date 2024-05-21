package service

import (
	"os"
	"testing"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
)

var stubNotificationURL = "https://example.com/source1/notification.json"
var stubSnapshot2URL = "https://example.com/ca128382-78d9-41d1-8927-1ecef15275be/nrtm-snapshot.2.047595d0fae972fbed0c51b4a41c7a349e0c47bb.json.gz"

func TestUpdateNRTMWithSourceInitialization(t *testing.T) {
	stubRepo := stubRepo{t: t, state: persist.NRTMState{}, err: &persist.ErrStateNotInitialized}
	stubClient := NewStubClient(t)
	tmpDir, err := os.MkdirTemp("", "nrtmtest*")
	if err != nil {
		t.Fatal("Could not create temp test directory")
	}
	defer func() {
		os.RemoveAll(tmpDir)
	}()
	UpdateNRTM(stubRepo, stubClient, stubNotificationURL, tmpDir)
}
