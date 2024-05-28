package service

import (
	"log"
	"os"
	"testing"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/pg"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/testresources"
)

var stubNotificationURL = "https://example.com/source1/notification.json"
var stubSnapshot2URL = "https://example.com/ca128382-78d9-41d1-8927-1ecef15275be/nrtm-snapshot.2.047595d0fae972fbed0c51b4a41c7a349e0c47bb.json.gz"

func TestE2EInstallSource(t *testing.T) {
	testresources.SetEnvVarsFromFile(t, "../testresources/env.test.conf")
	pgTestRepo := pgRepo()
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
	if err = processor.Connect(stubNotificationURL, ""); err != nil {
		t.Fatal("Failed to Connect", err)
	}
	sources, err := processor.ListSources()
	if len(sources) != 1 {
		t.Error("Should only be a single source")
	}
	src := sources[0]
	if src.Source != "EXAMPLE" {
		t.Error("Source should be EXAMPLE")
	}
	if src.Version != 3 {
		t.Error("Version should be 3")
	}
	if src.NotificationURL != stubNotificationURL {
		t.Error("NotificationURL should be", stubNotificationURL)
	}
	if src.SessionID != "XXX" {
		t.Error("SessionID should be XXX")
	}
}

func pgRepo() persist.Repository {
	dbURL := os.Getenv("PG_DATABASE_URL")
	if len(dbURL) == 0 {
		log.Fatal("ERROR no url for database", dbURL)
		return nil
	}
	repo := pg.PostgresRepository{}
	if err := repo.Initialize(dbURL); err != nil {
		log.Fatal("Failed to initialize repository")
	}
	return &repo
}
