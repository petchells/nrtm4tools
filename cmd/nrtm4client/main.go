package main

import (
	"log"
	"os"

	"github.com/petchells/nrtm4client/internal/nrtm4/cli"
	"github.com/petchells/nrtm4client/internal/nrtm4/service"
)

const mandatorySourceMessage = "Source name must be provided with the -source flag"

func main() {
	envVars := []string{"PG_DATABASE_URL", "NRTM4_FILE_PATH"}
	for _, ev := range envVars {
		if len(os.Getenv(ev)) <= 0 {
			log.Fatalln("Environment variable not set: ", ev)
		}
	}
	dbURL := os.Getenv("PG_DATABASE_URL")
	boltDBPath := os.Getenv("BOLT_DATABASE_PATH")
	nrtmFilePath := os.Getenv("NRTM4_FILE_PATH")
	config := service.AppConfig{
		NRTMFilePath:     nrtmFilePath,
		PgDatabaseURL:    dbURL,
		BoltDatabasePath: boltDBPath,
	}
	commander := cli.InitializeCommandProcessor(config)
	cli.Exec(commander)
}
