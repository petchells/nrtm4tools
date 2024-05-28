package persist

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/pg/db"
)

// NRTMFile is a binding to a PG database table
type NRTMFile struct {
	db.EntityManaged `em:"nrtm_file nf"`
	ID               uint64    `em:"."`
	Created          time.Time `em:"."`
	FileName         string    `em:"."`
	NRTMSourceID     uint64    `em:"."`
	Type             string    `em:"."`
	URL              string    `em:"."`
	Version          uint      `em:"."`
}

func GetLastState(tx pgx.Tx, source string) *NRTMFile {
	state := new(NRTMFile)
	descriptor := db.GetDescriptor(state)
	sql := fmt.Sprintf(`
		SELECT %v FROM %v
		WHERE
			st.source=$1
		ORDER BY
			st.version DESC,
			st.created DESC
		LIMIT 1
		`, strings.Join(descriptor.ColumnNamesWithAlias(), ", "), descriptor.TableNameWithAlias())
	rows, err := tx.Query(context.Background(), sql, source)
	if err != nil {
		log.Println("WARN", err)
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		log.Println(rows)
		err = rows.Scan(db.SelectValues(state)...)
		if err == nil {
			return state
		}
		log.Println("WARN scanning fields", err)
	}
	return nil
}
