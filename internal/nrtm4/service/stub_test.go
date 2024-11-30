package service

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"

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

func (r stubRepo) SaveSource(src persist.NRTMSource, notification persist.NotificationJSON) (persist.NRTMSource, error) {
	return persist.NRTMSource{}, nil
}

func (r stubRepo) Close() error {
	return nil
}

func (r stubRepo) SaveSnapshotObjects(source persist.NRTMSource, state persist.NRTMFile, rpslObject []rpsl.Rpsl) error {
	return nil
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

func (r stubRepo) AddModifyObject(src persist.NRTMSource, rpsl rpsl.Rpsl, file persist.NrtmFileJSON) error {
	return nil
}

func (r stubRepo) DeleteObject(src persist.NRTMSource, objectType string, primaryKey string, file persist.NrtmFileJSON) error {
	return nil
}

type stubClient struct {
	t *testing.T
}

func NewStubClient(t *testing.T) Client {
	return stubClient{t}
}

func (c stubClient) getUpdateNotification(url string) (persist.NotificationJSON, error) {
	var file persist.NotificationJSON
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
	} else if url == "https://example.com/ca128382-78d9-41d1-8927-1ecef15275be/nrtm-delta.3.d9c194acbb2cb0d4088c9d8a25d5871cdd802c79.json" {
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
	"version": 3,
	"snapshot": {
	  "version": 2,
	  "url": "https://example.com/ca128382-78d9-41d1-8927-1ecef15275be/nrtm-snapshot.2.047595d0fae972fbed0c51b4a41c7a349e0c47bb.json.gz",
	  "hash": "07396d52396a96b80eb4a5febbff2053ad945dbdbfdd020492d6fec7cf8cb526"
	},
	"deltas": [
	  {
		"version": 3,
		"url": "https://example.com/ca128382-78d9-41d1-8927-1ecef15275be/nrtm-delta.3.d9c194acbb2cb0d4088c9d8a25d5871cdd802c79.json",
		"hash": "bb65420644b598cdd7eb3b101f26ac033667d5edfe5c4f4fa005ff136e9eb8f8"
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
	"version": 2
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
