package testresources

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/petchells/nrtm4client/internal/nrtm4/pg/db"
)

// SetEnvVarsFromFile sets environment variables from a file
func SetEnvVarsFromFile(t *testing.T, fname string) {
	cnf, err := os.Open(fname)
	if err != nil {
		t.Fatal("Cannot open", fname, err)
	}
	defer cnf.Close()
	scanner := bufio.NewScanner(cnf)
	for scanner.Scan() {
		pair := strings.SplitN(scanner.Text(), "=", 2)
		os.Setenv(pair[0], pair[1])
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
}

// TruncateDatabase wipes all rows from all tables except '%schema_version' (Tern's version tracking table)
func TruncateDatabase(t *testing.T) {
	err := db.WithTransaction(func(tx pgx.Tx) error {
		selectSQL := `
			SELECT table_name
			FROM information_schema.tables
			WHERE table_schema='public'
				AND table_type='BASE TABLE'
				AND table_name not like '%schema_version'
				;
		`
		rows, err := tx.Query(context.Background(), selectSQL)
		if err != nil {
			return err
		}
		defer rows.Close()
		var tableNames []string
		for rows.Next() {
			var name string
			if err = rows.Scan(&name); err != nil {
				return err
			}
			tableNames = append(tableNames, name)
		}
		sql := fmt.Sprintf(`TRUNCATE %v CASCADE`, strings.Join(tableNames, ", "))
		_, err = tx.Exec(context.Background(), sql)
		return err
	})
	if err != nil {
		t.Fatal("Error when truncating tables in DB", err)
	}
}
