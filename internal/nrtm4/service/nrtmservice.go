package service

import (
	"errors"
	"log"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/nrtm4model"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
)

var ErrInvalidJSON = errors.New("invalid JSON")

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
			log.Printf("ERROR unknown action %v: '%v'\n", source, delta.Action)
			return ErrInvalidJSON
		}
	}
	return nil
}
