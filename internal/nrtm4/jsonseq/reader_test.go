package jsonseq

import (
	"encoding/json"
	"io"
	"testing"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/nrtm4model"
)

func TestJSONSequenceParser(t *testing.T) {
	// expectations
	sessionID := "ca128382-78d9-41d1-8927-1ecef15275be"
	numObjects := 3

	i := 0
	err := ReadStringRecords(snapshotExample, func(possJsonBytes []byte, err error) error {
		if i == 0 {
			snapshot := new(nrtm4model.SnapshotFile)
			err = json.Unmarshal(possJsonBytes, snapshot)
			if err != nil {
				t.Fatal(err)
			}
			if snapshot.SessionID != sessionID {
				t.Fatal("Expected", sessionID, "but was", snapshot.SessionID)
			}
		} else if i == 1 {
			object := new(nrtm4model.SnapshotObject)
			err = json.Unmarshal(possJsonBytes, object)
			if err != nil {
				t.Fatal(err)
			}
			if len(object.Object) < 10 {
				t.Fatal("Expected RPSL object string")
			}
		} else if i > 2 {
			t.Fatal("Expected three JSON entities")
		}
		i += 1
		return nil
	})
	if err != io.EOF {
		t.Fatal(err)
	}
	if i != numObjects {
		t.Fatal("Wrong number of JSON objects. Expected", numObjects, "but was", i)
	}
}

var snapshotExample = `
{
	"nrtm_version": 4,
	"type": "snapshot",
	"source": "EXAMPLE",
	"session_id": "ca128382-78d9-41d1-8927-1ecef15275be",
	"version": 3
}
{"object": "route: 192.0.2.0/24\norigin: AS65530\nsource: EXAMPLE"}
{"object": "route: 2001:db8::/32\norigin: AS65530\nsource: EXAMPLE"}
`
