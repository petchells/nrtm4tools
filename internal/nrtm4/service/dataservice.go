package service

import (
	"fmt"
	"log"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/nrtm4model"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
)

type ErrNRTMServiceError struct {
	Message string
}

func (e ErrNRTMServiceError) Error() string {
	return "nrtm service error: " + e.Message
}

func newNRTMServiceError(msg string, args ...any) ErrNRTMServiceError {
	return ErrNRTMServiceError{fmt.Sprintf(msg, args...)}
}

type NrtmDataService struct {
	Repository persist.Repository
}

func (ds NrtmDataService) ApplyDeltas(source string, deltas []nrtm4model.Change) error {
	for _, delta := range deltas {
		if delta.Action == "delete" {
			log.Println("i will delete", source, delta.PrimaryKey)
		} else if delta.Action == "add_modify" {
			log.Println("i will add/modify", source, delta.PrimaryKey)
		} else {
			return newNRTMServiceError("unknown delta action %v: '%v'", source, delta.Action)
		}
	}
	return nil
}
