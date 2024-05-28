package service

import "testing"

func TestLabelRegex(t *testing.T) {
	lbls := [...]string{
		"This_one_is-100.OK",
		"This_one_is not OK",
		"1not_ok",
		"ALSO-NOT!",
		"F",
	}
	for idx, lbl := range lbls {
		match := labelRegex.MatchString(lbl)
		if match && idx > 0 {
			t.Error("Invalid label should not match", lbl, match)
		} else if !match && idx == 0 {
			t.Error("Valid label should match", lbl, match)
		}
	}
}
