package nrtm4

import (
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/pg/db"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/pg/persist"
)

type AppConfig struct {
	NrtmUrlNotificationUrl string
	DatabaseURL            string
}

func Launch() {
	nrtmUrlNotificationUrl := os.Getenv("NRTM4_BASE_NOTIFICATION")
	dbUrl := os.Getenv("DATABASE_URL")
	repo := persist.PgRepository{}
	config := AppConfig{
		NrtmUrlNotificationUrl: nrtmUrlNotificationUrl,
		DatabaseURL:            dbUrl,
	}
	launchWithPg(repo, config)
}

func launchWithPg(repository persist.PgRepository, config AppConfig) {
	repository.InitializeConnectionPool(config.DatabaseURL)
	log.Println("DEBUG Launch()", config)
	err := db.WithTransaction(func(tx pgx.Tx) error {
		state := persist.GetLastState(tx)
		if state == nil {
			return ErrNoState
		}
		return nil
	})
	if err != nil {
		log.Println("ERROR", err)
	}

}
