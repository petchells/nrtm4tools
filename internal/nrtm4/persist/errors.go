package persist

import (
	"fmt"
)

var ErrStateNotInitialized = ErrNrtmClient{"state not initialized"}
var ErrNoEntity = ErrNrtmClient{"no json entity in record"}

type ErrNrtmClient struct {
	Msg string
}

func (e *ErrNrtmClient) Error() string {
	return fmt.Sprintf("ErrNrtmClient: %v", e.Msg)
}
