package nrtm4serve

import (
	"net/http"

	"github.com/petchells/nrtm4tools/internal/nrtm4/persist"
	"github.com/petchells/nrtm4tools/internal/nrtm4/service"
	"github.com/petchells/nrtm4tools/internal/nrtm4serve/rpc"
)

var (
	// Application codes from -32000 to -32098

	// Hash256ErrorCode -32010
	Hash256ErrorCode = -32010
	// NoDeltasInNotificationErrorCode -32020
	NoDeltasInNotificationErrorCode = -32020
	// SnapshotInsertFailedErrorCode -32030
	SnapshotInsertFailedErrorCode = -32030
	// DeltaUnavaliableErrorCode -32040
	DeltaUnavaliableErrorCode = -32040
	// NRTMServiceErrorCode -32050
	NRTMServiceErrorCode = -32050
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

// FetchSource returns a single source
func (api WebAPI) FetchSource(src, label string) *persist.NRTMSourceDetails {
	srcs, err := api.Processor.ListSources()
	if err != nil {
		logger.Warn("api.Processor.ListSources() returned error", "error", err)
		return nil
	}
	for _, s := range srcs {
		if s.Source == src && s.Label == label {
			return &s
		}
	}
	return nil
}

// ReplaceLabel replaces a label on a source
func (api WebAPI) ReplaceLabel(source, fromLabel, toLabel string) (*persist.NRTMSource, error) {
	src, err := api.Processor.ReplaceLabel(source, fromLabel, toLabel)
	return src, wrapErr(err)
}

// SaveProperties saves the properties for a source/label
func (api WebAPI) SaveProperties(source, label string, props persist.SourceProperties) (*persist.NRTMSource, error) {
	logger.Info("SaveProperties called", "source", source, "label", label, "props", props)
	src, err := api.Processor.SaveProperties(source, label, props)
	return src, err
}

// Connect connects a new source to the repo
func (api WebAPI) Connect(url, label string) (string, error) {
	err := api.Processor.Connect(url, label)
	if err != nil {
		service.UserLogger.Error("Connect failed", "url", url, "label", label, "error", err)
	}
	return "OK", wrapErr(err)
}

// Update updates a source to the latest version
func (api WebAPI) Update(src, label string) (persist.NRTMSourceDetails, error) {
	target, err := api.Processor.Update(src, label)
	if err != nil {
		return persist.NRTMSourceDetails{}, wrapErr(err)
	}
	deets, err := api.Processor.ListSources()
	if err != nil {
		return persist.NRTMSourceDetails{}, wrapErr(err)
	}
	for _, d := range deets {
		if d.ID == target.ID {
			return d, nil
		}
	}
	return persist.NRTMSourceDetails{}, nil
}

// RemoveSource removes a source from the repo
func (api WebAPI) RemoveSource(src, label string) (string, error) {
	return "OK", wrapErr(api.Processor.RemoveSource(src, label))
}

func wrapErr(err error) error {
	if err == nil {
		return nil
	}
	switch err {
	case service.ErrHashMismatch:
		return rpc.JSONRPCError{Code: Hash256ErrorCode, Message: err.Error()}
	case service.ErrNextConsecutiveDeltaUnavaliable:
		return rpc.JSONRPCError{Code: DeltaUnavaliableErrorCode, Message: err.Error()}
	case service.ErrSnapshotInsertFailed:
		return rpc.JSONRPCError{Code: SnapshotInsertFailedErrorCode, Message: err.Error()}
	case service.ErrSnapshotInsertFailed:
		return rpc.JSONRPCError{Code: SnapshotInsertFailedErrorCode, Message: err.Error()}
	case service.ErrNRTM4NoDeltasInNotification:
		return rpc.JSONRPCError{Code: NoDeltasInNotificationErrorCode, Message: err.Error()}
	}
	switch err.(type) {
	case service.ErrNRTMServiceError:
		return rpc.JSONRPCError{Code: NRTMServiceErrorCode, Message: err.Error()}
	}
	return err
}
