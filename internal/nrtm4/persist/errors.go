package persist

import (
	"fmt"
)

// ErrStateNotInitialized source is not initialized
var ErrStateNotInitialized = ErrNrtmClient{"source is not initialized"}

// ErrNoEntity expected to find some JSON
var ErrNoEntity = ErrNrtmClient{"no json found"}

// ErrNrtmClient something went wrong in the client
type ErrNrtmClient struct {
	Msg string
}

func (e *ErrNrtmClient) Error() string {
	return fmt.Sprintf("ErrNrtmClient: %v", e.Msg)
}
