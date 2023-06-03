package service

import (
	"log"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/nrtm4model"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
)

func UpdateNRTM(repo persist.Repository) {
	ds := NrtmDataService{Repository: repo}
	// Fetch notification
	// -- validate
	// -- new version?
	// Fetch state
	state, err := repo.GetState()
	if err != nil {
		log.Panicln("Failed to get state", err)
	}
	log.Println(state)
	// -- if no state, then initialize
	//    * get snapshot
	//    * save state
	//    * insert rpsl objects
	//    * done and dusted
	// -- compare with latest notification
	//    * is version newer? if not then blow
	//    * are there contiguous deltas since our last delta? if not, download snapshot
	//    * apply deltas
	ds.ApplyDeltas("RIR-TEST", []nrtm4model.Change{})
}
