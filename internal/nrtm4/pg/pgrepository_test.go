package pg

import (
	"strings"
	"testing"
	"unicode"
)

func TestSelectObjectSQL(t *testing.T) {
	sql := selectObjectQuery()

	expected := `
		SELECT rpsl.id rpsl_id, rpsl.object_type rpsl_object_type, rpsl.primary_key rpsl_primary_key, rpsl.nrtm_source_id rpsl_nrtm_source_id, rpsl.from_version rpsl_from_version, rpsl.to_version rpsl_to_version, rpsl.rpsl rpsl_rpsl
		FROM nrtm_rpslobject rpsl
		JOIN nrtm_source src ON src.id = rpsl.nrtm_source_id
		WHERE
			src.source ILIKE($1)
			AND UPPER(rpsl.primary_key) = UPPER($2)
			AND rpsl.object_type = UPPER($3)
			AND rpsl.to_version = 0`
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
