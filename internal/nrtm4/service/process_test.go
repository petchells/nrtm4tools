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
	defer func() {
		os.RemoveAll(tmpDir)
	}()
	conf := AppConfig{
		NRTMFilePath: tmpDir,
	}
	processor := NewNRTMProcessor(conf, pgTestRepo, stubClient)
	if err = processor.Connect(stubNotificationURL, ""); err != nil {
		t.Error("Failed to Connect", err)
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
