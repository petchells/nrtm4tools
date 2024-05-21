package service

import (
	"os"
	"strings"
	"testing"
)

var testResourcePath = "../test_resources/"

func TestGZIPSnapshotReader(t *testing.T) {
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
		counter++
		if err != nil {
			numErrors++
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

func TestPlainSnapshotReader(t *testing.T) {
	filename := "snapshot-sample.jsonseq"

	snapshotFile, err := os.Open(testResourcePath + filename)
	if err != nil {
		t.Fatal("Cannot open", filename)
	}
	t.Log("Opened", snapshotFile.Name())

	fm := fileManager{}
	numErrors := 0
	counter := 0
	fm.readSnapshotRecords(snapshotFile, func(bytes []byte, err error) error {
		counter++
		if err != nil {
			numErrors++
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

func TestWriteFromReaderToFile(t *testing.T) {
	file, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatal("Could not create temp file", err)
	}
	defer func() {
		file.Close()
		os.Remove(file.Name())
	}()

	fromStr := "Far and few, far and few are the lands where the Jumblies live."
	reader := strings.NewReader(fromStr)

	err = transferReaderToFile(reader, file)
	if err != nil {
		t.Error("Did not save file", err)
	}
}

func TestWriteFileToPath(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "nrtmtest")
	if err != nil {
		t.Fatal("Could not create temp directory", err)
	}
	defer func() {
		os.RemoveAll(tmpdir)
	}()

	fm := fileManager{
		client: NewStubClient(t),
	}

	file, err := fm.writeResourceToPath(stubSnapshot2URL, tmpdir)
	if file == nil || err != nil {
		t.Fatal("File was not written:", err)
	}

}
