package nrtm4serve

import (
	"net/http"

	"github.com/petchells/nrtm4tools/internal/nrtm4/persist"
	"github.com/petchells/nrtm4tools/internal/nrtm4/service"
	"github.com/petchells/nrtm4tools/internal/nrtm4serve/rpc"
)

var (
	// Application codes from -32000 to -32098

	// Hash256ErrorCode JSON RPC error code
	Hash256ErrorCode = -32010
	// SnapshotInsertFailedErrorCode JSON RPC error code
	SnapshotInsertFailedErrorCode = -32020
	// DeltaUnavaliableErrorCode JSON RPC error code
	DeltaUnavaliableErrorCode = -32030
	// NRTMServiceErrorCode problem with the service
	NRTMServiceErrorCode = -32040
)

// WebAPI defines the RPC functions used by the web client
type WebAPI struct {
	//	rpc.API
	Processor service.NRTMProcessor
}

// GetAuth implements rpc.API interface -- allows requests to all methods
func (api WebAPI) GetAuth(w http.ResponseWriter, r *http.Request, req rpc.JSONRPCRequest) (rpc.WebSession, bool) {
	return rpc.WebSession{}, true
}

// ListSources returns a list of sources
func (api WebAPI) ListSources() ([]persist.NRTMSourceDetails, error) {
	return api.Processor.ListSources()
}

// ReplaceLabel replaces a label on a source
func (api WebAPI) ReplaceLabel(source, fromLabel, toLabel string) (*persist.NRTMSource, error) {
	return api.Processor.ReplaceLabel(source, fromLabel, toLabel)
}

// Connect connects a new source to the repo
func (api WebAPI) Connect(url, label string) (string, error) {
	err := api.Processor.Connect(url, label)
	if err != nil {
		service.UserLogger.Error("Connect failed", "url", url, "label", label, "error", err)
	}
	return wrapResponse("OK", err)
}

// Update updates a source to the latest version
func (api WebAPI) Update(src, label string) (persist.NRTMSourceDetails, error) {
	target, err := api.Processor.Update(src, label)
	if err != nil {
		return wrapResponse(persist.NRTMSourceDetails{}, err)
	}
	deets, err := api.Processor.ListSources()
	if err != nil {
		return wrapResponse(persist.NRTMSourceDetails{}, err)
	}
	for _, d := range deets {
		if d.ID == target.ID {
			return wrapResponse(d, nil)
		}
	}
	return wrapResponse(persist.NRTMSourceDetails{}, nil)
}

// RemoveSource removes a source from the repo
func (api WebAPI) RemoveSource(src, label string) (string, error) {
	return wrapResponse("OK", api.Processor.RemoveSource(src, label))
}

func wrapResponse[T any](res T, err error) (T, error) {
	if err == nil {
		return res, nil
	}
	if err == service.ErrHashMismatch {
		return res, rpc.JSONRPCError{Code: Hash256ErrorCode, Message: err.Error()}
	} else if err == service.ErrNextConsecutiveDeltaUnavaliable {
		return res, rpc.JSONRPCError{Code: DeltaUnavaliableErrorCode, Message: err.Error()}
	} else if err == service.ErrSnapshotInsertFailed {
		return res, rpc.JSONRPCError{Code: SnapshotInsertFailedErrorCode, Message: err.Error()}
	}
	switch err.(type) {
	case service.ErrNRTMServiceError:
		return res, rpc.JSONRPCError{Code: NRTMServiceErrorCode, Message: err.Error()}
	}
	return res, err
}
