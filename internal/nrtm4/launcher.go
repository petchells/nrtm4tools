package nrtm4

import (
	"log"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/pg"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/service"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/util"
)

var logger = util.Logger

// Connect sets up the execution environment then invokes the connect command
func Connect(config service.AppConfig, url, label string) {
	var httpClient service.HTTPClient
	repo := pg.PostgresRepository{}
	if err := repo.Initialize(config.PgDatabaseURL); err != nil {
		log.Fatal("Failed to initialize repository")
	}
	defer repo.Close()
	processor := service.NewNRTMProcessor(config, repo, httpClient)
	commander := service.NewCommandProcessor(processor)
	commander.Connect(url, label)
}

// Update sets up the execution environment then invokes the update command
func Update(config service.AppConfig, sourceName, label string) {
	var httpClient service.HTTPClient
	repo := pg.PostgresRepository{}
	if err := repo.Initialize(config.PgDatabaseURL); err != nil {
		log.Fatal("Failed to initialize repository")
	}
	defer repo.Close()
	processor := service.NewNRTMProcessor(config, repo, httpClient)
	commander := service.NewCommandProcessor(processor)
	commander.Update(sourceName, label)
}
