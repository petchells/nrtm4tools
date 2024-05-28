package nrtm4

import (
	"log"
	"os"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/pg"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/service"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/util"
)

var logger = util.Logger

// LaunchPg launch with PostgreSQL database
func LaunchPg(config service.AppConfig) {
	repo := pg.PostgresRepository{}
	if err := repo.Initialize(config.PgDatabaseURL); err != nil {
		log.Fatal("Failed to initialize repository")
	}
	defer repo.Close()
	connect(&repo, config)
}

// LaunchBolt launch with Bolt database
// func LaunchBolt(config AppConfig) {
// 	repo := bolt.BoltRepository{}
// 	if err := repo.Initialize(config.BoltDatabasePath); err != nil {
// 		log.Fatal("Failed to initialize repository")
// 	}
// 	defer repo.Close()
// 	update(&repo, config)
// }

func connect(repo persist.Repository, config service.AppConfig) {
	logger.Debug("Launch()", "config", config)
	logger.Debug("Arguments", "args", os.Args)
	if len(os.Args) < 3 {
		logger.Error("Not sure what to do. Exiting.")
		return
	}
	var httpClient service.HTTPClient
	commander := service.NewCommandProcessor(config, repo, httpClient)
	commander.Connect(os.Args[2], "")
}
