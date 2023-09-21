package nrtm4

import (
	"log"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/bolt"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/pg"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/service"
)

type AppConfig struct {
	NrtmUrlNotificationUrl string
	PgDatabaseURL          string
	NrtmFilePath           string
	BoltDatabasePath       string
}

func Launch(config AppConfig) {
	//repo := bolt.BoltRepository{}
	repo := pg.PgRepository{}
	launchWithPg(&repo, config)
}

func launchWithPg(repository *pg.PgRepository, config AppConfig) {
	repository.Initialize(config.PgDatabaseURL)
	log.Println("DEBUG Launch()", config)
	var httpClient service.HttpClient
	service.UpdateNRTM(repository, httpClient, config.NrtmUrlNotificationUrl, config.NrtmFilePath)
}

func launchWithBolt(repository *bolt.BoltRepository, config AppConfig) {
	repository.Initialize(config.BoltDatabasePath)
	log.Println("DEBUG Launch()", config)
	var httpClient service.HttpClient
	service.UpdateNRTM(repository, httpClient, config.NrtmUrlNotificationUrl, config.NrtmFilePath)
	repository.Close()
}
