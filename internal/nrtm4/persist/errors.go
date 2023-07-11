package persist

import (
	"fmt"
)

var ErrNoState = ErrNrtmClient{"state not initialized"}
var ErrNoEntity = ErrNrtmClient{"no json entity in record"}
var ErrFetchingState = ErrNrtmClient{"no state exists"}

type ErrNrtmClient struct {
	Msg string
}

func (e *ErrNrtmClient) Error() string {
	return fmt.Sprintf("ErrNrtmClient: %v", e.Msg)
}
