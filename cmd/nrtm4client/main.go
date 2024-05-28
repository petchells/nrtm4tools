package main

import (
	"log"
	"os"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/service"
)

func main() {
	envVars := []string{"PG_DATABASE_URL", "NRTM4_FILE_PATH", "BOLT_DATABASE_PATH"}
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
	// TODO
	// Parse multiple URLs and file path
	// Start one goroutine for each source
	nrtm4.LaunchPg(config)
}
