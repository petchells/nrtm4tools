package persist

import (
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/pg/db"
)

type PgRepository struct {
	persist.Repository
}

func (repo PgRepository) InitializeConnectionPool(dbUrl string) {
	db.InitializeConnectionPool(dbUrl)
}
