package main

import (
	"log"
	"os"

	"gitlab.com/etchells/nrtm4-client/internal/nrtm4"
)

func main() {
	envVars := []string{"DATABASE_URL", "NRTM4_BASE_NOTIFICATION"}
	for _, ev := range envVars {
		if len(os.Getenv(ev)) <= 0 {
			log.Fatalln("Environment variable not set: ", ev)
		}
	}
	nrtm4.Launch()
}
