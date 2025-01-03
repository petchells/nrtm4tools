package service

import (
	"log"
	"os"
	"sort"
	"testing"

	"github.com/petchells/nrtm4client/internal/nrtm4/persist"
	"github.com/petchells/nrtm4client/internal/nrtm4/pg"
	"github.com/petchells/nrtm4client/internal/nrtm4/testresources"
)

var stubNotificationURL = "https://example.com/source1/notification.json"
var stubSnapshot2URL = "https://example.com/ca128382-78d9-41d1-8927-1ecef15275be/nrtm-snapshot.2.047595d0fae972fbed0c51b4a41c7a349e0c47bb.json.gz"

func TestFileRefSorter(t *testing.T) {
	refs := []persist.FileRefJSON{
		{
			Version: 4,
			URL:     "https://xxx.xxx.xxx/4",
			Hash:    "4444",
		},
		{
			Version: 6,
			URL:     "https://xxx.xxx.xxx/6",
			Hash:    "6666",
		},
		{
			Version: 3,
			URL:     "https://xxx.xxx.xxx/3",
			Hash:    "3333",
		},
		{
			Version: 5,
			URL:     "https://xxx.xxx.xxx/5",
			Hash:    "5",
		},
	}
	sort.Sort(fileRefsByVersion(refs))
	expect := [...]uint32{3, 4, 5, 6}
	for idx, v := range expect {
		if refs[idx].Version != v {
			t.Error("Expected", v, "but got", refs[idx].Version)
		}
	}
}

func TestE2EConnect(t *testing.T) {
	testresources.SetEnvVarsFromFile(t, "../testresources/env.test.conf")
	pgTestRepo := pgRepo()
	testresources.TruncateDatabase(t)
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
	if src.SessionID != "ca128382-78d9-41d1-8927-1ecef15275be" {
		t.Error("SessionID should be", "ca128382-78d9-41d1-8927-1ecef15275be")
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
