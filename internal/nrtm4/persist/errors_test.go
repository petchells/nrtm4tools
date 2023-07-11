package persist

import "testing"

func TestErrors(t *testing.T) {
	myErr1 := ErrNrtmClient{"oo"}
	myErr2 := ErrNrtmClient{"oo"}
	myErr3 := ErrNrtmClient{"ee"}
	if myErr1 == myErr2 {
		t.Log("They are equal")
	} else {
		t.Error("They are not equal, but they should be")
	}
	if myErr1 == myErr3 {
		t.Error("They should *not* be the same")
	} else {
		t.Log("They are not equal")
	}

}
