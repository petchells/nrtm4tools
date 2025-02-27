package persist

import (
	"testing"

	"github.com/petchells/nrtm4tools/internal/nrtm4/pg/db"
)

func TestNRTMSourceMapping(t *testing.T) {
	t.Log("Amazing")
}

func TestColumnNameConversionFromFieldTags(t *testing.T) {
	expected := [...]string{"id", "source", "session_id", "version", "notification_url", "label", "created"}
	o := NRTMSource{}
	dtor := db.GetDescriptor(&o)
	names := dtor.ColumnNames()
	if len(expected) != len(names) {
		t.Errorf("Expected '%d' fields but got '%d'", len(expected), len(names))
	}
	for i := 0; i < len(expected); i++ {
		if expected[i] != names[i] {
			t.Errorf("Expected field name '%v' but got '%v'", expected[i], names[i])
		}
	}
}
