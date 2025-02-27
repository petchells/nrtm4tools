package main

import (
	"flag"
	"log"
	"os"

	"github.com/petchells/nrtm4tools/internal/nrtm4/service"
	"github.com/petchells/nrtm4tools/internal/nrtm4serve"
)

var port = flag.Int("port", 8080, "server port number")
var webdir = flag.String("webdir", "", "path to static web root")

func main() {
	flag.Parse()
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
	nrtm4serve.Launch(config, *port, *webdir)
}
