package nrtm4serve

import (
	"net/http"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/service"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4serve/rpc"
)

// WebAPI defines the RPC functions used by the web client
type WebAPI struct {
	rpc.API
	Repo      persist.Repository
	AppConfig service.AppConfig
}

// GetAuth implements interface -- allows requests to all methods
func (api WebAPI) GetAuth(w http.ResponseWriter, r *http.Request, req rpc.JSONRPCRequest) (rpc.WebSession, bool) {
	return rpc.WebSession{}, true
}

// ListSources returns a list of sources
func (api WebAPI) ListSources() ([]persist.NRTMSourceDetails, error) {
	var httpClient service.HTTPClient
	processor := service.NewNRTMProcessor(api.AppConfig, api.Repo, httpClient)
	return processor.ListSources()
}
