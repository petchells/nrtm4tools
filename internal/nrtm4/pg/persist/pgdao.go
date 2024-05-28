package persist

import "github.com/jackc/pgx/v5"

// PgDao dao for postgresql
type PgDao struct {
	Tx pgx.Tx
}
