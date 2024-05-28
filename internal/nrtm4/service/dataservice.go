package service

import (
	"fmt"
	"log"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
)

// ErrNRTMServiceError is when sth is wrong with the NRTM server
type ErrNRTMServiceError struct {
	Message string
}

func (e ErrNRTMServiceError) Error() string {
	return "nrtm service error: " + e.Message
}

func newNRTMServiceError(msg string, args ...any) ErrNRTMServiceError {
	return ErrNRTMServiceError{fmt.Sprintf(msg, args...)}
}

// NrtmDataService is an implementation of a persist.Repository
type NrtmDataService struct {
	Repository persist.Repository
}

// func (ds NrtmDataService) applyDeltas(source string, deltas []nrtm4model.DeltaJSON) error {
// 	for _, delta := range deltas {
// 		if delta.Action == nrtm4model.DeltaDeleteAction {
// 			log.Println("i will delete", source, delta.PrimaryKey)
// 		} else if delta.Action == nrtm4model.DeltaAddModifyAction {
// 			log.Println("i will add/modify", source, delta.PrimaryKey)
// 		} else {
// 			return newNRTMServiceError("unknown delta action %v: '%v'", source, delta.Action)
// 		}
// 	}
// 	return nil
// }

func (ds NrtmDataService) getSourceByURLAndLabel(url string, label string) *persist.NRTMSource {
	sources, err := ds.getSources()
	if err != nil {
		log.Panicln("Failure calling NrtmDataService.getSources", err)
	}
	for _, src := range sources {
		if src.NotificationURL == url && src.Label == label {
			found := src
			return &found
		}
	}
	return nil
}

func (ds NrtmDataService) getSources() ([]persist.NRTMSource, error) {
	return ds.Repository.GetSources()
}

func (ds NrtmDataService) saveNewSource(source persist.NRTMSource) (persist.NRTMSource, error) {
	return ds.Repository.SaveSource(source)
}
