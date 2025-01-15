package pg

import (
	"strings"
	"testing"
	"unicode"
)

func TestSelectObjectSQL(t *testing.T) {
	sql := selectCurrentObjectQuery()

	expected := `
		SELECT id, object_type, primary_key, nrtm_source_id, from_version, to_version, rpsl
		FROM nrtm_rpslobject
		WHERE
			nrtm_source_id = $1
			AND primary_key = UPPER($2)
			AND object_type = UPPER($3)
			AND to_version = 0`
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
