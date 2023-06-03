package service

import (
	"fmt"
	"log"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/nrtm4model"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
)

type ErrInvalidJSON struct {
	Message string
}

func (e ErrInvalidJSON) Error() string {
	return "invalid JSON: " + e.Message
}

func newInvalidJSONError(msg string, args ...any) ErrInvalidJSON {
	return ErrInvalidJSON{fmt.Sprintf(msg, args...)}
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
			return newInvalidJSONError("unknown action %v: '%v'", source, delta.Action)
		}
	}
	return nil
}
