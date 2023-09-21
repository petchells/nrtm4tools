package service

import (
	"os"
	"testing"
)

var testResourcePath = "../test_resources/"

func TestSnapshotReader(t *testing.T) {
	filename := "snapshot-sample.jsonseq.gz"

	snapshotFile, err := os.Open(testResourcePath + filename)
	if err != nil {
		t.Fatal("Cannot open", filename)
	}
	t.Log("Opened", snapshotFile.Name())

	fm := fileManager{}
	numErrors := 0
	counter := 0
	fm.readSnapshotRecords(snapshotFile, func(bytes []byte, err error) error {
		counter += 1
		if err != nil {
			numErrors += 1
			return err
		}
		return nil
	})
	if numErrors != 1 {
		t.Error("Expected only one (EOF) error, but was", numErrors)
	}
	if counter != 10 {
		t.Error("Expected to read 10 lines, but was", counter)
	}
}
