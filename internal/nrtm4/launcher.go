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
	repo := pg.PgRepository{}
	repo.Initialize(config.PgDatabaseURL)
	update(&repo, config)
}

func LaunchBolt(config AppConfig) {
	repo := bolt.BoltRepository{}
	repo.Initialize(config.BoltDatabasePath)
	update(&repo, config)
}

func update(repo persist.Repository, config AppConfig) error {
	log.Println("DEBUG Launch()", config)
	var httpClient service.HttpClient
	service.UpdateNRTM(repo, httpClient, config.NrtmUrlNotificationUrl, config.NrtmFilePath)
	return repo.Close()
}
