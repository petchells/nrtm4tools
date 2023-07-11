package service

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/nrtm4model"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
)

func TestUpdateNRTMWithSourceInitialization(t *testing.T) {
	stubRepo := stubRepo{t: t, state: persist.NRTMState{}, err: &persist.ErrNoState}
	stubClient := stubClient{t}
	tmpDir := filepath.Join(os.TempDir(), "/nrtmtest")
	defer func() {
		os.RemoveAll(tmpDir)
	}()
	UpdateNRTM(stubRepo, stubClient, "https://example.com/source1/notification.json", tmpDir)
}

type stubRepo struct {
	t     *testing.T
	state persist.NRTMState
	err   *persist.ErrNrtmClient
}

func (r stubRepo) InitializeConnectionPool(dbUrl string) {}

func (r stubRepo) GetState(source string) (persist.NRTMState, *persist.ErrNrtmClient) {
	state := r.state
	if r.err != nil {
		log.Println("ERROR GetState", r.err)
		return state, &persist.ErrFetchingState
	}
	if source == "EXAMPLE" {
		return state, nil
	}
	r.t.Fatal("Unexpected request for source", source)
	return state, &persist.ErrFetchingState
}

func (r stubRepo) SaveState(state persist.NRTMState) error {
	expected := "notification.json"
	if state.FileName == expected {
		return nil
	}
	r.t.Error("SaveState failed. expected file name", expected, "but was", state.FileName)
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

func (c stubClient) getResponseBody(url string) (io.Reader, error) {
	var reader io.Reader
	if url == "https://example.com/ca128382-78d9-41d1-8927-1ecef15275be/nrtm-snapshot.2.047595d0fae972fbed0c51b4a41c7a349e0c47bb.json.gz" {
		reader = strings.NewReader(snapshotExample)
		var buf bytes.Buffer
		zw := gzip.NewWriter(&buf)
		zw.Write([]byte(snapshotExample))
		return gzip.NewReader(&buf)
	} else if url == "https://example.com/ca128382-78d9-41d1-8927-1ecef15275be/nrtm-delta.1.784a2a65aba22e001fd25a1b9e8544e058fbc703.json" {
		reader = strings.NewReader(deltaExample)
	} else {
		c.t.Error("Call to unexpected URL", url)
		return reader, errors.New("unexpected file url")
	}
	return reader, nil
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

var snapshotExample = `
{
	"nrtm_version": 4,
	"type": "snapshot",
	"source": "EXAMPLE",
	"session_id": "ca128382-78d9-41d1-8927-1ecef15275be",
	"version": 3
}
{"object": "route: 192.0.2.0/24\norigin: AS65530\nsource: EXAMPLE"}
{"object": "route: 2001:db8::/32\norigin: AS65530\nsource: EXAMPLE"}
`

var deltaExample = `
{
	"nrtm_version": 4,
	"type": "delta",
	"source": "EXAMPLE",
	"session_id": "ca128382-78d9-41d1-8927-1ecef15275be",
	"version": 3
}
{
	"action": "delete",
	"object_class": "person",
	"primary_key": "PRSN1-EXAMPLE"
}
{
	"action": "delete",
	"object_class": "route",
	"primary_key": "192.0.2.0/24AS65530"
}
{
	"action": "add_modify",
	"object": "route: 2001:db8::/32\norigin: AS65530\nsource: EXAMPLE"
}
`
