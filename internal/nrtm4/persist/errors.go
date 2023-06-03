package persist

import (
	"errors"
	"fmt"
)

var ErrNoState = errors.New("state not initialized")

type ErrOne struct {
	Msg string
}

func (e ErrOne) Error() string {
	return fmt.Sprintf("everthing is wrong: %v", e.Msg)
}
