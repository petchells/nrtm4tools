package persist

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/pg/db"
)

type NRTMState struct {
	db.EntityManaged `em:"nrtmstate st"`
	ID               uint64    `em:"."`
	Source           string    `em:"."`
	Version          uint      `em:"."`
	URL              string    `em:"."`
	Type             string    `em:"."`
	FileName         string    `em:"."`
	Created          time.Time `em:"."`
}

func GetLastState(tx pgx.Tx, source string) *NRTMState {
	state := new(NRTMState)
	descriptor := db.GetDescriptor(state)
	sql := fmt.Sprintf(`
		SELECT %v FROM %v
		WHERE
			source=$1
		  AND 
			st.version=MAX(st.version)
		`, descriptor.ColumnNamesWithAlias(), descriptor.TableNameWithAlias())
	log.Println(sql)
	rows, err := tx.Query(context.Background(), sql, source)
	if err != nil {
		log.Println("WARN", err)
		return nil
	}
	defer rows.Close()
	var states []NRTMState
	for rows.Next() {
		log.Println(rows)
		err = rows.Scan(db.FieldValues(state)...)
		if err != nil {
			log.Println("WARN scanning fields", err)
			return state
		}
		states = append(states, *state)
	}
	log.Println("DEBUG states", states)

	// if err != nil {
	// 	log.Println("WARN", err)
	// }
	return state
}
