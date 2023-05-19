package nrtm4

import (
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"gitlab.com/etchells/nrtm4-client/internal/nrtm4/db"
	"gitlab.com/etchells/nrtm4-client/internal/nrtm4/persist"
)

func Launch() {
	db.InitializeConnectionPool(os.Getenv("DATABASE_URL"))
	nrtmUrlNotificationUrl := os.Getenv("NRTM4_BASE_NOTIFICATION")
	log.Println("DEBUG Launch()", nrtmUrlNotificationUrl)
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
