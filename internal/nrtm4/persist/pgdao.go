package persist

import "github.com/jackc/pgx/v5"

type PgDao struct {
	Tx pgx.Tx
}
