package service

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/nrtm4model"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/rpsl"
)

type stubRepo struct {
	t     *testing.T
	state persist.NRTMFile
	err   error
}

func (r stubRepo) Initialize(dbURL string) error {
	return nil
}

func (r stubRepo) GetSources() []persist.NRTMSource {
	return []persist.NRTMSource{}
}

func (r stubRepo) CreateSource(label string, source string, notificationURL string, pathOnDisk string) (*persist.NRTMSource, error) {
	return &persist.NRTMSource{}, nil
}

func (r stubRepo) Close() error {
	return nil
}

func (r stubRepo) SaveSnapshotFile(source persist.NRTMSource, tate persist.NRTMFile, snapshotFile nrtm4model.SnapshotFileJSON) error {
	return nil
}

func (r stubRepo) SaveSnapshotObject(source persist.NRTMSource, state persist.NRTMFile, rpslObject rpsl.Rpsl) error {
	return nil
}

func (r stubRepo) SaveSnapshotObjects(source persist.NRTMSource, state persist.NRTMFile, rpslObject []rpsl.Rpsl) error {
	return nil
}

func (r stubRepo) GetState(source string) (persist.NRTMFile, error) {
	state := r.state
	if r.err != nil {
		return state, &persist.ErrStateNotInitialized
	}
	if source == "EXAMPLE" {
		return state, nil
	}
	r.t.Fatal("Unexpected request for source", source)
	return state, &persist.ErrStateNotInitialized
}

func (r stubRepo) SaveFile(nrtmFile *persist.NRTMFile) error {
	expected := "notification.json"
	if nrtmFile.FileName == expected {
		r.state = *nrtmFile
		return nil
	}
	r.t.Error("SaveFile failed. expected file name", expected, "but was", nrtmFile.FileName)
	return nil
}

type stubClient struct {
	t *testing.T
}

func NewStubClient(t *testing.T) Client {
	return stubClient{t}
}

func (c stubClient) getUpdateNotification(url string) (nrtm4model.NotificationJSON, error) {
	var file nrtm4model.NotificationJSON
	if url == stubNotificationURL {
		json.Unmarshal([]byte(notificationExample), &file)
		return file, nil
	}
	c.t.Error("unexpected notification url")
	return file, errors.New("unexpected notification url")
}

func (c stubClient) getResponseBody(url string) (io.Reader, error) {
	var reader io.Reader
	if url == stubSnapshot2URL {
		reader = strings.NewReader(snapshotExample)
		var buf bytes.Buffer
		zw := gzip.NewWriter(&buf)
		zw.Write([]byte(snapshotExample))
		zw.Close()
		return bytes.NewReader(buf.Bytes()), nil
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
