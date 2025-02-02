package service

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/petchells/nrtm4client/internal/nrtm4/persist"
	"github.com/petchells/nrtm4client/internal/nrtm4/util"
)

var testResourcePath = "../testresources/"

func TestSuccess(t *testing.T) {
	fm := fileManager{dlClientStub{}}
	_, errs := fm.downloadNotificationFile("")
	if len(errs) > 0 {
		t.Error("should not be any errors but found:", errs[0])
	} else {
		t.Log("OK")
	}
}

type dlClientStub struct {
	Client
}

func (c dlClientStub) getUpdateNotification(string) (persist.NotificationJSON, error) {
	notification := persist.NotificationJSON{
		NrtmFileJSON: persist.NrtmFileJSON{
			NrtmVersion: 4,
			SessionID:   uuid.NewString(),
			Type:        persist.NotificationFile.String(),
			Source:      "ZZZZ",
			Version:     22,
		},
		Timestamp:      util.AppClock.Now().Format(time.RFC3339),
		NextSigningKey: new(string),
		SnapshotRef: persist.FileRefJSON{
			URL: "https://xxx.xxx.xx/notification.json",
		},
		DeltaRefs: &[]persist.FileRefJSON{},
	}
	return notification, nil
}

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

	snapshotFile, err := os.Open(testResourcePath + filename)
	if err != nil {
		t.Fatal("Cannot open", filename)
	}
	t.Log("Opened", snapshotFile.Name())

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

func TestValidateURLString(t *testing.T) {
	type testURL struct {
		str      string
		expected bool
	}
	testURLs := []testURL{
		{"https://nrtm4.example.com/nrtm4/notification.json", true},
		{"ftp://nrtm4.example.com/nrtm4/notification.json", false},
		{"RIPE/nrtm-snapshot.374234.RIPE.db44e038-1f07-4d54-a307-1b32339f141a.7755dc0a05b5024dd092a7a68d1b7b0.json.gz", false},
	}
	for _, turl := range testURLs {
		if validateURLString(turl.str) != turl.expected {
			t.Error("Validation failed. Expected", turl.expected, "for:", turl.str, validateURLString(turl.str))
		}
	}
}

func TestFetchFileAndCheckHash(t *testing.T) {
	body := `{
	"secretMessage", "Some text in a file"
	}
	`
	ref := persist.FileRefJSON{
		URL:     "https://nrtmv4.example.com/testtext.txt",
		Hash:    "123456",
		Version: 3,
	}
	client := stubDeltaClient{
		responseBody: body,
	}
	dir, err := os.MkdirTemp("", "nrtm4test")
	if err != nil {
		t.Fatal("Failed to create temp dir", err)
	}
	defer os.RemoveAll(dir)
	{
		fm := fileManager{
			client: client,
		}
		_, err := fm.fetchFileAndCheckHash(ref, dir)
		if err != ErrHashMismatch {
			t.Fatal("Expected ErrHashMismatch but was:", err)
		}
	}
	{
		fm := fileManager{
			client: client,
		}
		ref.Hash = "4d14d44910c1abae9b55b6cc0f722369834b3c1942f3ee4bc0e051b1de10794d"
		f, err := fm.fetchFileAndCheckHash(ref, dir)
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
