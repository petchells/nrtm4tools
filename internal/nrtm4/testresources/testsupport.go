package testresources

import (
	"bufio"
	"os"
	"strings"
	"testing"
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
