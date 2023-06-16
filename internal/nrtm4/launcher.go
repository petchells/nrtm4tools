package nrtm4

import (
	"log"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/pg/persist"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/service"
)

type AppConfig struct {
	NrtmUrlNotificationUrl string
	PgDatabaseURL          string
	NrtmFilePath           string
}

func Launch(config AppConfig) {
	repo := persist.PgRepository{}
	launchWithPg(repo, config)
}

func launchWithPg(repository persist.PgRepository, config AppConfig) {
	repository.InitializeConnectionPool(config.PgDatabaseURL)
	log.Println("DEBUG Launch()", config)
	var httpClient service.HttpClient
	service.UpdateNRTM(repository, httpClient, config.NrtmUrlNotificationUrl, config.NrtmFilePath)
}
