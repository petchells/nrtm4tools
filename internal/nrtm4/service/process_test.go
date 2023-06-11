package service

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/nrtm4model"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
)

type stubRepo struct {
	t *testing.T
}

func (r stubRepo) InitializeConnectionPool(dbUrl string) {}

func (r stubRepo) GetState(source string) (persist.NRTMState, error) {
	var state persist.NRTMState
	if source == "EXAMPLE" {
		return state, nil
	}
	r.t.Fatal("Unexpected request for source", source)
	return state, errors.New("unknown source")
}

func (r stubRepo) SaveState(state persist.NRTMState) error {
	expected := persist.NRTMState{}
	if expected != state {
		r.t.Error("SaveState failed. expected", expected, "but was", state)
	}
	return nil
}

type stubClient struct {
	t *testing.T
}

func (c stubClient) getUpdateNotification(url string) (nrtm4model.Notification, error) {
	var file nrtm4model.Notification
	if url == "https://example.com/source1/notification.json" {
		json.Unmarshal([]byte(notificationExample), &file)
		return file, nil
	}
	c.t.Error("unexpected notification url")
	return file, errors.New("unexpected notification url")
}

func (c stubClient) fetchFile(url string) (io.ReadCloser, error) {
	var reader io.ReadCloser
	file := "{}"
	if url == "https://example.com/ca128382-78d9-41d1-8927-1ecef15275be/nrtm-snapshot.2.047595d0fae972fbed0c51b4a41c7a349e0c47bb.json.gz" {
		reader.Read([]byte(file))
	} else if url == "https://example.com/ca128382-78d9-41d1-8927-1ecef15275be/nrtm-delta.1.784a2a65aba22e001fd25a1b9e8544e058fbc703.json" {
		reader.Read([]byte(file))
	} else {
		c.t.Error("Call to unexpected URL", url)
		return reader, errors.New("unexpected file url")
	}
	return reader, nil
}

func TestUpdateNRTM(t *testing.T) {
	stubClient := stubClient{t}
	stubRepo := stubRepo{t}
	tmpDir := filepath.Join(os.TempDir(), "/nrtmtest")
	defer func() {
		os.RemoveAll(tmpDir)
	}()
	UpdateNRTM(stubRepo, stubClient, "https://example.com/source1/notification.json", tmpDir)
}

var notificationExample = `
{
	"nrtm_version": 4,
	"timestamp": "2022-01-00T15:00:00Z",
	"type": "notification",
	"next_signing_key": "96..ae",
	"source": "EXAMPLE",
	"session_id": "ca128382-78d9-41d1-8927-1ecef15275be",
	"version": 4,
	"snapshot": {
	  "version": 3,
	  "url": "https://example.com/ca128382-78d9-41d1-8927-1ecef15275be/nrtm-snapshot.2.047595d0fae972fbed0c51b4a41c7a349e0c47bb.json.gz",
	  "hash": "9a..86"
	},
	"deltas": [
	  {
		"version": 2,
		"url": "https://example.com/ca128382-78d9-41d1-8927-1ecef15275be/nrtm-delta.1.784a2a65aba22e001fd25a1b9e8544e058fbc703.json",
		"hash": "62..a2"
	  },
	  {
		"version": 3,
		"url": "https://example.com/ca128382-78d9-41d1-8927-1ecef15275be/nrtm-delta.2.0f681f07cfab5611f3681bf030ec9f6fa3442fb0.json",
		"hash": "25..9a"
	  },
	  {
		"version": 4,
		"url": "https://example.com/ca128382-78d9-41d1-8927-1ecef15275be/nrtm-delta.3.d9c194acbb2cb0d4088c9d8a25d5871cdd802c79.json",
		"hash": "b4..13"
	  }
	]
  }  
`
