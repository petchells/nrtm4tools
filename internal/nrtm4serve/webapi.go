package nrtm4serve

import (
	"net/http"

	"github.com/petchells/nrtm4client/internal/nrtm4/persist"
	"github.com/petchells/nrtm4client/internal/nrtm4/service"
	"github.com/petchells/nrtm4client/internal/nrtm4serve/rpc"
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

// RemoveSource removes a source from the repo
func (api WebAPI) RemoveSource(src, label string) (string, error) {
	err := api.Processor.RemoveSource(src, label)
	return "", err
}
