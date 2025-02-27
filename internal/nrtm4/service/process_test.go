package service

import (
	"log"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/petchells/nrtm4tools/internal/nrtm4/persist"
	"github.com/petchells/nrtm4tools/internal/nrtm4/testresources"
)

const (
	baseURL             = "https://example.com/source1/"
	stubNotificationURL = "notification.json"
)

type labelExpectation struct {
	label  string
	expect bool
}

func TestLabelRegex(t *testing.T) {
	lbls := [...]labelExpectation{{
		"This_one_is-100.OK", true},
		{"1_is_ok", true},
		{"YES$nowerky", false},
		{"No\\BS", false},
		{"-1", true},
		{"F", true},
		{"1970-01-01", true},
		{"This one is OK", true},
		{"    This one will be trimmed   ", true},
		{"-------", false},
		{"------1", true},
	}
	for _, lbl := range lbls {
		match := labelRe.MatchString(lbl.label)
		if match != lbl.expect {
			if lbl.expect {
				t.Error("Label regex should succeed", lbl.label)
			} else {
				t.Error("Label regex should fail", lbl.label)
			}
		}
	}
}

func TestFileRefSorter(t *testing.T) {
	refs := []persist.FileRefJSON{
		{
			Version: 4,
			URL:     "https://xxx.xxx.xxx/4",
			Hash:    "4444",
		},
		{
			Version: 6,
			URL:     "https://xxx.xxx.xxx/6",
			Hash:    "6666",
		},
		{
			Version: 3,
			URL:     "https://xxx.xxx.xxx/3",
			Hash:    "3333",
		},
		{
			Version: 5,
			URL:     "https://xxx.xxx.xxx/5",
			Hash:    "5",
		},
	}
	sort.Sort(fileRefsByVersion(refs))
	expect := [...]int64{3, 4, 5, 6}
	for idx, v := range expect {
		if refs[idx].Version != v {
			t.Error("Expected", v, "but got", refs[idx].Version)
		}
	}
}

func TestFindUpdatesSuccess(t *testing.T) {

	var notification persist.NotificationJSON
	testresources.ReadTestJSONToPtr(t, "ripe-notification-file.json", &notification)
	source := stubsource()

	fileRefs, err := findUpdates(notification, source)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
	expectedLen := 9
	if len(fileRefs) != expectedLen {
		t.Fatalf("Unexpected slice length. Expected %d but was %d", expectedLen, len(fileRefs))
	}
}

func TestValidateNotificationErrors(t *testing.T) {

	var notification persist.NotificationJSON
	{
		testresources.ReadTestJSONToPtr(t, "ripe-notification-file.json", &notification)
		refs := notification.DeltaRefs
		dr := append(refs[:10], refs[11:]...)
		notification.DeltaRefs = dr

		expect := ErrNRTM4NotificationDeltaSequenceBroken

		err := validateNotificationFile(notification)

		if err != expect {
			t.Errorf("Expected error %v but was %v", expect, err)
		}
	}
	{
		testresources.ReadTestJSONToPtr(t, "ripe-notification-file.json", &notification)
		refs := notification.DeltaRefs
		dr := refs[:len(refs)-2]
		notification.DeltaRefs = dr

		expect := ErrNRTM4NotificationVersionDoesNotMatchDelta

		err := validateNotificationFile(notification)
		if err != expect {
			t.Errorf("Expected error %v but was %v", expect, err)
		}
	}
	{
		testresources.ReadTestJSONToPtr(t, "ripe-notification-file.json", &notification)
		refs := notification.DeltaRefs
		dr := append(refs[:10], refs[9:]...)
		notification.DeltaRefs = dr

		expect := ErrNRTM4DuplicateDeltaVersion

		err := validateNotificationFile(notification)
		if err != expect {
			t.Errorf("Expected error %v but was %v", expect, err)
		}
	}
	{
		testresources.ReadTestJSONToPtr(t, "ripe-notification-file.json", &notification)
		dr := []persist.FileRefJSON{}
		notification.DeltaRefs = dr

		expect := ErrNRTM4NoDeltasInNotification

		err := validateNotificationFile(notification)
		if err != expect {
			t.Errorf("Expected error %v but was %v", expect, err)
		}
	}
}

func TestFullURLFunction(t *testing.T) {
	base := "https://nrtm.example.eu/path/to/nrtm4/notification-file.json"
	{
		rel := "source/r2d1-65535.EXAMPLE.json"
		expected := "https://nrtm.example.eu/path/to/nrtm4/source/r2d1-65535.EXAMPLE.json"

		result := fullURL(base, rel)

		if result != expected {
			t.Error("fullURL returned wrong url expected:", expected, "but was:", result)
		}
	}
	{
		rel := "/source/r2d1-65535.EXAMPLE.json"
		expected := "https://nrtm.example.eu/path/to/nrtm4/source/r2d1-65535.EXAMPLE.json"

		result := fullURL(base, rel)

		if result != expected {
			t.Error("fullURL returned wrong url expected:", expected, "but was:", result)
		}
	}
	{
		base = "nrtm.example.eu"
		rel := "/source/r2d1-65535.EXAMPLE.json"
		expected := ""

		result := fullURL(base, rel)

		if result != expected {
			t.Error("fullURL returned wrong url expected:", expected, "but was:", result)
		}
	}
}

func TestValidateURLString(t *testing.T) {
	type testURL struct {
		str      string
		expected bool
	}
	testURLs := []testURL{
		{"https://nrtm4.example.com/nrtm4/notification.json", true},
		{"https:///nrtm4.example.com/nrtm4/notification.json", true},
		{"ftp://nrtm4.example.com/nrtm4/notification.json", false},
		{"RIPE/nrtm-snapshot.374234.RIPE.db44e038-1f07-4d54-a307-1b32339f141a.7755dc0a05b5024dd092a7a68d1b7b0.json.gz", false},
	}
	for _, turl := range testURLs {
		if validateURLString(turl.str) != turl.expected {
			t.Error("Validation failed. Expected", turl.expected, "for:", turl.str, validateURLString(turl.str))
		}
	}
}

func TestConnectErrors(t *testing.T) {
	// Set up
	pgTestRepo := mockRepo{}
	stubClient := NewTestClient(t, baseURL, "version2to6", "unf_2-4.json")
	tmpDir, err := os.MkdirTemp("", "nrtmtest*")
	if err != nil {
		t.Fatal("Could not create temp test directory")
	}
	defer os.RemoveAll(tmpDir)

	conf := AppConfig{
		NRTMFilePath: tmpDir,
	}
	processor := NewNRTMProcessor(conf, pgTestRepo, stubClient)

	// Run test
	label := "what a load of testing"
	{
		if err = processor.Connect("not a url", label); err != ErrBadNotificationURL {
			t.Error("Bad URL should fail", err)
		}
	}
	{
		if err = processor.Connect(baseURL+stubNotificationURL, "-=-"); err != ErrInvalidLabel {
			t.Error("Bad label should fail", err)
		}
	}
	{
		str := "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"
		str += str
		if err = processor.Connect(baseURL+stubNotificationURL, str); err != ErrInvalidLabel {
			t.Error("Bad label should fail", err)
		}
	}
	{
		sources := []persist.NRTMSource{
			{
				ID:              987,
				NotificationURL: baseURL + stubNotificationURL,
				Label:           label,
			},
		}
		pgTestRepo.sources = sources
		processor.repo = pgTestRepo
		if err = processor.Connect(baseURL+stubNotificationURL, label); err != ErrSourceAlreadyExists {
			t.Error("Source already exist, should be rejected", err)
		}
	}
}

func stubsource() persist.NRTMSource {
	t, err := time.Parse(time.RFC3339, "2025-01-04T23:01:00Z")
	if err != nil {
		log.Fatalln("bad timestamp")
	}
	src := persist.NRTMSource{
		ID:              576576257634,
		Source:          "TEST_SRC",
		SessionID:       "db44e038-1f07-4d54-a307-1b32339f141a",
		Version:         350684,
		NotificationURL: baseURL + stubNotificationURL,
		Label:           "",
		Created:         t,
	}
	return src
}
