package testresources

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/petchells/nrtm4client/internal/nrtm4/pg/db"
)

// SetEnvVarsFromEnvTestFile creates env vars from env.test.conf
func SetEnvVarsFromEnvTestFile(t *testing.T) {
	base := pathToPackage()
	SetEnvVarsFromFile(t, base+"/env.test.conf")
}

// SetEnvVarsFromFile creates env vars from the given file
func SetEnvVarsFromFile(t *testing.T, fname string) {
	cnf := openFile(t, fname)
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

// OpenFile returns an open file
func OpenFile(t *testing.T, fname string) *os.File {
	return openFile(t, filepath.Join(pathToPackage(), fname))
}

// ReadJSON reads a JSON file relative to this dir and unmarshalls it to a pointer
func ReadJSON(t *testing.T, jsonFile string, ptr any) {
	jsonPath := filepath.Join(pathToPackage(), jsonFile)
	readJSON(t, jsonPath, ptr)
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

func readJSON(t *testing.T, fileName string, ptr any) error {
	var err error

	file := openFile(t, fileName)
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, ptr)
}

func openFile(t *testing.T, fname string) *os.File {
	var err error
	var file *os.File
	if file, err = os.Open(fname); err != nil {
		log.Println(err)
		t.Fatal("File does not exist", fname)
		return nil
	}
	return file
}

func pathToPackage() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Println("Test cannot determine the path to package: testsupport")
	}
	return filepath.Dir(filename)
}
