package nrtm4

import (
	"log"

	"github.com/petchells/nrtm4client/internal/nrtm4/pg"
	"github.com/petchells/nrtm4client/internal/nrtm4/service"
	"github.com/petchells/nrtm4client/internal/nrtm4/util"
)

var logger = util.Logger

// InitializeCommandProcessor starts a db connection pool
func InitializeCommandProcessor(config service.AppConfig) service.CommandExecutor {
	var httpClient service.HTTPClient
	repo := pg.PostgresRepository{}
	if err := repo.Initialize(config.PgDatabaseURL); err != nil {
		log.Fatal("Failed to initialize repository")
	}
	defer repo.Close()
	processor := service.NewNRTMProcessor(config, repo, httpClient)
	return service.NewCommandProcessor(processor)
}
