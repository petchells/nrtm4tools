package service

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/petchells/nrtm4client/internal/nrtm4/persist"
	"github.com/petchells/nrtm4client/internal/nrtm4/testresources"
)

func TestConnectWithPgRepo(t *testing.T) {
	// Set up
	pgTestRepo := testresources.SetTestEnvAndInitializePG(t)
	stubClient := NewStubClient(t)
	tmpDir, err := os.MkdirTemp("", "nrtmtest*")
	if err != nil {
		t.Fatal("Could not create temp test directory")
	}
	defer os.RemoveAll(tmpDir)

	conf := AppConfig{
		NRTMFilePath: tmpDir,
	}
	processor := NewNRTMProcessor(conf, pgTestRepo, stubClient)

	// Run test
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
		if s.Source == "EXAMPLE" && s.Label == label {
			src = s
			break
		}
	}
	if src.Source != "EXAMPLE" {
		t.Error("Source should be EXAMPLE")
	}
	if src.Version != 3 {
		t.Error("Version should be 3")
	}
	if src.NotificationURL != baseURL+stubNotificationURL {
		t.Error("NotificationURL should be", baseURL+stubNotificationURL)
	}
	if src.SessionID != "ca128382-78d9-41d1-8927-1ecef15275be" {
		t.Error("SessionID should be", "ca128382-78d9-41d1-8927-1ecef15275be")
	}

	err = processor.Update("example", label)

	if err != nil {
		t.Error("Error update returned an error", err)
	}
}
