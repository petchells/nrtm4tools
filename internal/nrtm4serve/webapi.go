package nrtm4serve

import (
	"net/http"

	"github.com/petchells/nrtm4tools/internal/nrtm4/persist"
	"github.com/petchells/nrtm4tools/internal/nrtm4/service"
	"github.com/petchells/nrtm4tools/internal/nrtm4serve/rpc"
)

var (
	DeltaUnavaliableErrCode = -32020
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
	return "OK", err
}

// Update updates a source to the latest version
func (api WebAPI) Update(src, label string) (string, error) {
	err := api.Processor.Update(src, label)
	if err == service.ErrNextConsecutiveDeltaUnavaliable {
		return "", rpc.JSONRPCError{Code: DeltaUnavaliableErrCode, Message: err.Error()}
	}
	return "OK", err
}

// RemoveSource removes a source from the repo
func (api WebAPI) RemoveSource(src, label string) (string, error) {
	err := api.Processor.RemoveSource(src, label)
	return "OK", err
}
