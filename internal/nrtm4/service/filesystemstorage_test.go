package service

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/petchells/nrtm4client/internal/nrtm4/persist"
	"github.com/petchells/nrtm4client/internal/nrtm4/testresources"
)

func TestSuccess(t *testing.T) {
	fm := fileManager{NewTestClient(t, baseURL, "version2to6", "unf_2-4.json")}
	_, err := fm.downloadNotificationFile("")
	if err != nil {
		t.Error("should not be any errors but found:", err)
	} else {
		t.Log("OK")
	}
}

func TestGZIPSnapshotReader(t *testing.T) {
	filename := "snapshot-sample.jsonseq.gz"
	snapshotFile := testresources.OpenFile(t, filename)

	fm := fileManager{}
	numErrors := 0
	counter := 0
	fm.readJSONSeqRecords(snapshotFile, func(bytes []byte, err error) error {
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
	snapshotFile := testresources.OpenFile(t, filename)

	fm := fileManager{}
	numErrors := 0
	counter := 0
	err := fm.readJSONSeqRecords(snapshotFile, func(bytes []byte, err error) error {
		counter++
		if err != nil {
			numErrors++
			return err
		}
		return nil
	})
	if numErrors != 1 {
		t.Error("Expected only one error, but was", numErrors)
	}
	if err != io.EOF {
		t.Error("Expected EOF but was", err)
	}
	if counter != 10 {
		t.Error("Expected to read 10 lines, but was", counter)
	}
}

func TestWriteFileToPath(t *testing.T) {
	// Given...
	tmpdir, err := os.MkdirTemp("", "nrtmtest*")
	if err != nil {
		t.Fatal("Could not create temp directory", err)
	}
	defer func() {
		os.RemoveAll(tmpdir)
	}()

	fm := fileManager{
		client: NewTestClient(t, baseURL, "version2to6", "unf_2-4.json"),
	}
	snapshotPath := "snapshot.2.TEST.jsonseq.gz"

	// When...
	file, err := fm.writeResourceToPath(baseURL+snapshotPath, filepath.Join(tmpdir, snapshotPath))
	if file == nil || err != nil {
		t.Fatal("File was not written:", err)
	}

	// Then...
	expected := filepath.Join(tmpdir, snapshotPath)
	if file.Name() != expected {
		t.Error("File name failed. Expected", expected, "but was", file.Name())
	}

}

func TestFetchFileAndCheckHash(t *testing.T) {
	unfURL := "https://wherever.eu/unf.json"
	body := `{
	"secretMessage", "Some text in a file"
	}
	`
	ref := persist.FileRefJSON{
		URL:     "testtext.txt",
		Hash:    "123456",
		Version: 3,
	}
	client := stubDeltaClient{
		responseBody: body,
	}
	dir, err := os.MkdirTemp("", "nrtm4test*")
	if err != nil {
		t.Fatal("Failed to create temp dir", err)
	}
	defer os.RemoveAll(dir)
	{
		fm := fileManager{
			client: client,
		}
		_, err := fm.fetchFileAndCheckHash(unfURL, ref, dir)
		if err != ErrHashMismatch {
			t.Fatal("Expected ErrHashMismatch but was:", err)
		}
	}
	{
		fm := fileManager{
			client: client,
		}
		ref.Hash = "4d14d44910c1abae9b55b6cc0f722369834b3c1942f3ee4bc0e051b1de10794d"
		f, err := fm.fetchFileAndCheckHash(unfURL, ref, dir)
		if err != nil {
			t.Fatal("Unexpected error:", err)
		}
		defer f.Close()
		bytes, err := os.ReadFile(f.Name())
		if err != nil {
			t.Fatal("Could not read file:", f.Name(), err)
		}
		if string(bytes) != body {
			t.Error("Got unexpected body", string(bytes))
		}
	}

}
