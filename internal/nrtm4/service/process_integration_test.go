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

func TestConnectWithPgRepo(t *testing.T) {
	// Set up
	tmpDir, err := os.MkdirTemp("", "nrtmtest*")
	if err != nil {
		t.Fatal("Could not create temp test directory")
	}
	defer os.RemoveAll(tmpDir)
	conf := AppConfig{
		NRTMFilePath: tmpDir,
	}
	pgTestRepo := testresources.SetTestEnvAndInitializePG(t)
	stubClient := NewTestClient(t, baseURL, "version2to6", "unf_2-4.json")
	processor := NewNRTMProcessor(conf, pgTestRepo, stubClient)

	// Run test
	srcname := "TEST"
	label := filepath.Base(tmpDir)
	if err = processor.Connect(baseURL+stubNotificationURL, label); err != nil {
		t.Fatal("Failed to Connect", err)
	}

	// Assertions
	sources, err := processor.ListSources()
	if len(sources) < 1 {
		t.Error("Should be at least one source")
	}
	var src persist.NRTMSourceDetails
	for _, s := range sources {
		if s.Source == srcname && s.Label == label {
			src = s
			break
		}
	}
	if src.Source != srcname {
		t.Error("Source should be", srcname)
	}
	if src.Version != 4 {
		t.Error("Version should be 4")
	}
	if src.NotificationURL != baseURL+stubNotificationURL {
		t.Error("NotificationURL should be", baseURL+stubNotificationURL)
	}
	if src.SessionID != "17db6715-18ae-410f-973e-47981b52f023" {
		t.Error("SessionID should be", "17db6715-18ae-410f-973e-47981b52f023")
	}

	stubClient.NotificationFileName("unf2_6.json")
	err = processor.Update(strings.ToLower(srcname), label)

	if err != nil {
		t.Error("Error update returned an error", err)
	}
}

type tcConfig struct {
	baseURL, testDataDir, notifile string
}

type TestClient struct {
	conf tcConfig
	t    *testing.T
}

func NewTestClient(t *testing.T, baseURL, testDataDir, notifile string) TestClient {
	return TestClient{
		conf: tcConfig{baseURL, testDataDir, notifile},
		t:    t,
	}
}

func (c TestClient) NotificationFileName(fname string) {
	c.conf.notifile = fname
}

func (c TestClient) getUpdateNotification(_ string) (persist.NotificationJSON, error) {
	var notifile persist.NotificationJSON
	fname := filepath.Join(c.conf.testDataDir, c.conf.notifile)
	testresources.ReadTestJSONToPtr(c.t, fname, &notifile)
	return notifile, nil
}

func (c TestClient) getResponseBody(requrl string) (io.Reader, error) {
	if !strings.HasPrefix(requrl, c.conf.baseURL) {
		c.t.Fatal("Request for unrecognizer URL", requrl)
	}
	fname := requrl[len(c.conf.baseURL):]
	fpath := filepath.Join(c.conf.testDataDir, fname)
	return testresources.OpenFile(c.t, fpath), nil
}
