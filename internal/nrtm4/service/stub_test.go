package service

import (
	"testing"

	"github.com/petchells/nrtm4tools/internal/nrtm4/persist"
	"github.com/petchells/nrtm4tools/internal/nrtm4/rpsl"
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
	// 			SourceID: src.ID,
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
