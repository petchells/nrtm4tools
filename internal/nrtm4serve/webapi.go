package nrtm4serve

import (
	"net/http"

	"github.com/petchells/nrtm4tools/internal/nrtm4/persist"
	"github.com/petchells/nrtm4tools/internal/nrtm4/service"
	"github.com/petchells/nrtm4tools/internal/nrtm4serve/rpc"
)

var (
	// DeltaUnavaliableErrCode JSON RPC error code
	DeltaUnavaliableErrCode = -32060
	// Hash256ErrCode JSON RPC error code
	Hash256ErrCode = -32020
	// SnapshotInsertFailedErrCode JSON RPC error code
	SnapshotInsertFailedErrCode = -32040
)

// WebAPI defines the RPC functions used by the web client
type WebAPI struct {
	//	rpc.API
	Processor service.NRTMProcessor
}

// GetAuth implements interface -- allows requests to all methods
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
func (api WebAPI) Update(src, label string) (string, error) {
	_, err := api.Processor.Update(src, label)
	return wrapResponse("OK", err)
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
		return res, rpc.JSONRPCError{Code: Hash256ErrCode, Message: err.Error()}
	} else if err == service.ErrNextConsecutiveDeltaUnavaliable {
		return res, rpc.JSONRPCError{Code: DeltaUnavaliableErrCode, Message: err.Error()}
	} else if err == service.ErrSnapshotInsertFailed {
		return res, rpc.JSONRPCError{Code: SnapshotInsertFailedErrCode, Message: err.Error()}
	}
	return res, err
}
