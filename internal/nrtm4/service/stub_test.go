package service

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/petchells/nrtm4client/internal/nrtm4/persist"
	"github.com/petchells/nrtm4client/internal/nrtm4/rpsl"
)

type stubRepo struct {
	t     *testing.T
	state persist.NRTMFile
	err   error
}

func (r *stubRepo) Initialize(dbURL string) error {
	return nil
}

func (r stubRepo) ListSources() []persist.NRTMSource {
	return []persist.NRTMSource{}
}

func (r *stubRepo) SaveSource(src persist.NRTMSource, notification persist.NotificationJSON) (persist.NRTMSource, error) {
	return persist.NRTMSource{}, nil
}

func (r *stubRepo) Close() error {
	return nil
}

func (r *stubRepo) SaveSnapshotObjects(source persist.NRTMSource, state persist.NRTMFile, rpslObject []rpsl.Rpsl) error {
	return nil
}

func (r *stubRepo) SaveFile(nrtmFile *persist.NRTMFile) error {
	expected := "notification.json"
	if nrtmFile.FileName == expected {
		r.state = *nrtmFile
		return nil
	}
	r.t.Error("SaveFile failed. expected file name", expected, "but was", nrtmFile.FileName)
	return nil
}

func (r *stubRepo) AddModifyObject(src persist.NRTMSource, rpsl rpsl.Rpsl, file persist.NrtmFileJSON) error {
	return nil
}

func (r *stubRepo) DeleteObject(src persist.NRTMSource, objectType string, primaryKey string, file persist.NrtmFileJSON) error {
	return nil
}

type stubClient struct {
	t *testing.T
}

type mockRepo struct {
	persist.Repository
	sources []persist.NRTMSource
}

func (mr mockRepo) SaveSource(source persist.NRTMSource, notifile persist.NotificationJSON) (persist.NRTMSource, error) {
	id := uint64((len(mr.sources) + 1000))
	src := source
	src.ID = id
	// deets := persist.NRTMSourceDetails{
	// 	NRTMSource: src,
	// 	Notifications: []persist.Notification{
	// 		{
	// 			ID:           id,
	// 			Version:      uint32(notifile.Version),
	// 			NRTMSourceID: src.ID,
	// 			Payload:      notifile,
	// 			Created:      util.AppClock.Now(),
	// 		},
	// 	},
	// }
	mr.sources = append(mr.sources, src)
	return src, nil
}

func (mr mockRepo) ListSources() ([]persist.NRTMSource, error) {
	return mr.sources, nil
}

func NewStubClient(t *testing.T) Client {
	return stubClient{t}
}

func (c stubClient) getUpdateNotification(url string) (persist.NotificationJSON, error) {
	var file persist.NotificationJSON
	if url == baseURL+stubNotificationURL {
		json.Unmarshal([]byte(notificationExample), &file)
		return file, nil
	}
	c.t.Error("unexpected notification url")
	return file, errors.New("unexpected notification url")
}

func (c stubClient) getResponseBody(requrl string) (io.Reader, error) {
	var reader io.Reader
	if requrl == baseURL+stubSnapshot2URL {
		reader = strings.NewReader(snapshotExample)
		var buf bytes.Buffer
		zw := gzip.NewWriter(&buf)
		zw.Write([]byte(snapshotExample))
		zw.Close()
		return bytes.NewReader(buf.Bytes()), nil
	} else if requrl == baseURL+delta3URL {
		reader = strings.NewReader(deltaExample)
	} else {
		c.t.Error("Call to unexpected URL", requrl)
		return reader, errors.New("unexpected file url")
	}
	return reader, nil
}

var notificationExample = fmt.Sprintf(`
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
	  "url": "%v",
	  "hash": "7e1a5a4763ddb399feae52e45036ce7218877e7ca2c6dd20fec82efeb03074c0"
	},
	"deltas": [
	  {
		"version": 3,
		"url": "%v",
		"hash": "d829e802200605c79a0c42331a7a1333b524b05110c321b1bda6a390fbeaf6e2"
	  }
	]
  }
`, stubSnapshot2URL, delta3URL)

var snapshotExample = `
{
	"nrtm_version": 4,
	"type": "snapshot",
	"source": "EXAMPLE",
	"session_id": "ca128382-78d9-41d1-8927-1ecef15275be",
	"version": 2
}
{"object": "route: 192.0.2.0/24\norigin: AS65530\nsource: EXAMPLE"}
{"object": "route6: 2001:db8::/32\norigin: AS65530\nsource: EXAMPLE"}
{"object": "person: Bob the Builder\nnic-hdl: PRSN1-EXAMPLE\nsource: EXAMPLE"}
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
	"object": "route6:        2001:db8::/32\norigin:        AS65530\nsource: EXAMPLE"
}
`
