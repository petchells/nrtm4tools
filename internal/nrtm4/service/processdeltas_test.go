package service

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/petchells/nrtm4client/internal/nrtm4/persist"
	"github.com/petchells/nrtm4client/internal/nrtm4/testresources"
)

func TestApplyDeltas(t *testing.T) {

	label := "TestApplyDeltas"

	var err error
	repo := testresources.SetTestEnvAndInitializePG(t)
	f := testresources.OpenFile(t, "nrtm-delta.multiple-ops-same-pk.jsonseq")
	defer f.Close()
	bytes, _ := io.ReadAll(f)

	tmpDir, err := os.MkdirTemp("", "nrtm4*")
	if err != nil {
		t.Fatal("Could not create temp dir")
	}
	defer os.RemoveAll(tmpDir)

	config := AppConfig{
		NRTMFilePath: tmpDir,
	}
	client := stubDeltaClient{
		responseBody: string(bytes),
	}

	p := NRTMProcessor{
		repo:   repo,
		config: config,
		client: client,
	}
	deltas := []persist.FileRefJSON{
		{
			URL:     "n3.json",
			Version: 3,
			Hash:    "6e938ff1642485a651bf7cf14cd31c44eca17515909d8ddd9ed01efc840a61b1",
		},
	}
	notification := persist.NotificationJSON{
		NrtmFileJSON: persist.NrtmFileJSON{
			NrtmVersion: 4,
			Version:     3,
		},
		DeltaRefs: deltas,
	}
	source := persist.NRTMSource{
		Version:         2,
		SessionID:       "db44e038-1f07-4d54-a307-1b32339f141a",
		Source:          "RIPE",
		NotificationURL: "http://test.test.test/unf.json",
		Label:           label,
	}
	fdir := filepath.Join(tmpDir, source.Source, source.SessionID)
	if err = os.MkdirAll(fdir, 0775); err != nil {
		t.Fatal("Failed to create directory for files")
	}
	srcs, _ := repo.ListSources()
	for _, src := range srcs {
		if src.Source == source.Source && src.Label == source.Label {
			repo.RemoveSource(src)
		}
	}
	if source, err = repo.SaveSource(source, notification); err != nil {
		t.Fatal("Failed to save source")
	}

	err = syncDeltas(p, notification, source)

	if err != nil {
		t.Error("Failed to apply deltas", err)
	}
}

type stubDeltaClient struct {
	notification persist.NotificationJSON
	responseBody string
}

func (c stubDeltaClient) getUpdateNotification(string) (persist.NotificationJSON, error) {
	return c.notification, nil
}

func (c stubDeltaClient) getResponseBody(string) (io.Reader, error) {
	rdr := strings.NewReader(c.responseBody)
	return rdr, nil
}
