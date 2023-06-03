package nrtm4

import (
	"log"
	"os"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/pg/persist"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/service"
)

type AppConfig struct {
	NrtmUrlNotificationUrl string
	DatabaseURL            string
	SnapshotPath           string
}

func Launch() {
	nrtmUrlNotificationUrl := os.Getenv("NRTM4_BASE_NOTIFICATION")
	dbUrl := os.Getenv("DATABASE_URL")
	snapshotPath := os.Getenv("SNAPSHOT_PATH")
	repo := persist.PgRepository{}
	config := AppConfig{
		NrtmUrlNotificationUrl: nrtmUrlNotificationUrl,
		DatabaseURL:            dbUrl,
		SnapshotPath:           snapshotPath,
	}
	launchWithPg(repo, config)
}

func launchWithPg(repository persist.PgRepository, config AppConfig) {
	repository.InitializeConnectionPool(config.DatabaseURL)
	log.Println("DEBUG Launch()", config)
	service.UpdateNRTM(repository)
}
