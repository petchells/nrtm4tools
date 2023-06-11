package main

import (
	"log"
	"os"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4"
)

func main() {
	envVars := []string{"DATABASE_URL", "NRTM4_NOTIFICATION_URL", "NRTM4_FILE_PATH"}
	for _, ev := range envVars {
		if len(os.Getenv(ev)) <= 0 {
			log.Fatalln("Environment variable not set: ", ev)
		}
	}
	nrtmUrlNotificationUrl := os.Getenv("NRTM4_BASE_NOTIFICATION")
	dbUrl := os.Getenv("DATABASE_URL")
	nrtmFilePath := os.Getenv("NRTM4_FILE_PATH")
	config := nrtm4.AppConfig{
		NrtmUrlNotificationUrl: nrtmUrlNotificationUrl,
		NrtmFilePath:           nrtmFilePath,
		DatabaseURL:            dbUrl,
	}
	nrtm4.Launch(config)
}
