package nrtm4

import (
	"log"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/pg"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/service"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/util"
)

var logger = util.Logger

// LaunchPg launch with PostgreSQL database
func LaunchPg(config service.AppConfig, args []string) {
	repo := pg.PostgresRepository{}
	if err := repo.Initialize(config.PgDatabaseURL); err != nil {
		log.Fatal("Failed to initialize repository")
	}
	defer repo.Close()
	launch(&repo, config, args)
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

func launch(repo persist.Repository, config service.AppConfig, args []string) {
	logger.Debug("Launch()", "config", config)
	logger.Info("Arguments", "args", args)
	if len(args) < 2 {
		logger.Error("Not sure what to do. Exiting.")
		return
	}
	var httpClient service.HTTPClient
	processor := service.NewNRTMProcessor(config, repo, httpClient)
	commander := service.NewCommandProcessor(processor)
	cmd := args[0]
	if cmd == "connect" {
		url := args[1]
		var label string
		if len(args) > 3 {
			logger.Error("Not a command")
		}
		if len(args) == 3 {
			label = args[2]
		}
		commander.Connect(url, label)
	} else if cmd == "update" {
		var sourceName = args[1]
		var label string
		if len(args) > 3 {
			logger.Error("Not a command")
		}
		if len(args) == 3 {
			label = args[2]
		}
		commander.Update(sourceName, label)
	} else {
		logger.Error("Not a command")
	}
}
