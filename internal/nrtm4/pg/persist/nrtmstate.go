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
	Created          time.Time `em:"."`
	URL              string    `em:"."`
	FileVersion      uint      `em:"."`
	IsDelta          bool      `em:"."`
	Delta            string    `em:"."`
	SnapshotPath     string    `em:"."`
	Source           uint64    `em:"."`
}

func GetLastState(tx pgx.Tx) *NRTMState {
	state := new(NRTMState)
	descriptor := db.GetDescriptor(state)
	sql := fmt.Sprintf(`SELECT %v FROM %v WHERE st.version=MAX(st.version)`, descriptor.ColumnNamesWithAlias(), descriptor.TableNameWithAlias())
	log.Println(sql)
	err := tx.QueryRow(context.Background(), sql).Scan(db.FieldValues(state)...)
	if err != nil {
		log.Println("WARN", err)
	}
	return state
}
