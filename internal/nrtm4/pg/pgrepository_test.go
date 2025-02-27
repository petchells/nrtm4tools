package pg

import (
	"strings"
	"testing"
	"unicode"
)

func TestGetSources(t *testing.T) {
	// repo := testresources.SetTestEnvAndInitializePG(t)
	// sources, err := repo.ListSources()

	// if err != nil {
	// 	t.Error("Should not be an error when getting souces")
	// }
	// if len(sources) != 0 {
	// 	t.Error("Should not be any sources. Unless other tests are running.")
	// }
}

func TestSelectObjectSQL(t *testing.T) {
	sql := selectCurrentObjectQuery()

	expected := `
		SELECT id, object_type, primary_key, nrtm_source_id, version, rpsl
		FROM nrtm_rpslobject
		WHERE
			nrtm_source_id = $1
			AND primary_key = UPPER($2)
			AND object_type = UPPER($3)`
	if reduceWhiteSpace(sql) != reduceWhiteSpace(expected) {
		t.Errorf("Got unexpected SQL\n%v\nbut wanted\n%v\n", sql, expected)
	}
}

func TestReduceWhiteSpace(t *testing.T) {
	input := [...]string{
		"How now     brown      cow",
		"    How now   \n\n  brown   \t\r   cow\n  ",
	}
	expected := [...]string{
		"How now brown cow",
		"How now brown cow",
	}
	for idx, str := range input {
		got := reduceWhiteSpace(str)
		if expected[idx] != got {
			t.Error("WS reducer doesn't work")
		}
	}
}

func reduceWhiteSpace(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	var lastWasWS = true
	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
			lastWasWS = false
		} else {
			if !lastWasWS {
				b.WriteRune(' ')
				lastWasWS = true
			}
		}
	}
	return strings.TrimRight(b.String(), " ")
}
