package nrtm4

import (
	"log"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/bolt"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/pg"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/service"
)

type AppConfig struct {
	NrtmUrlNotificationUrl string
	PgDatabaseURL          string
	NrtmFilePath           string
	BoltDatabasePath       string
}

func LaunchPg(config AppConfig) {
	repo := pg.PostgresRepository{}
	if err := repo.Initialize(config.PgDatabaseURL); err != nil {
		log.Fatal("Failed to initialize repository")
	}
	defer repo.Close()
	update(&repo, config)
}

func LaunchBolt(config AppConfig) {
	repo := bolt.BoltRepository{}
	if err := repo.Initialize(config.BoltDatabasePath); err != nil {
		log.Fatal("Failed to initialize repository")
	}
	defer repo.Close()
	update(&repo, config)
}

func update(repo persist.Repository, config AppConfig) {
	log.Println("DEBUG Launch()", config)
	var httpClient service.HttpClient
	service.UpdateNRTM(repo, httpClient, config.NrtmUrlNotificationUrl, config.NrtmFilePath)
}
