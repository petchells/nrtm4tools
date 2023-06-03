package service

import (
	"log"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/nrtm4model"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
)

func UpdateNRTM(repo persist.Repository, url string) {
	// Fetch notification
	// -- validate
	// -- new version?
	var notification nrtm4model.Notification
	var err error

	if notification, err = getUpdateNotification(url); err != nil {
		log.Println("ERROR failed to fetch notificationFile", err)
		return
	}
	if errs := validateNotificationFile(notification); len(errs) > 0 {
		for _, err := range errs {
			log.Println("ERROR notificationFile failed validation", err)
		}
		return
	}

	// Fetch state
	ds := NrtmDataService{Repository: repo}
	state, err := repo.GetState(notification.Source)
	if err != nil {
		log.Println("Failed to get state", err)
		//repo.CreateState(state)
		// -- if no state, then initialize
		//    * get snapshot
		getSnapshot(notification)
		//    * save state
		//    * insert rpsl objects
		//    * done and dusted
		return
	}
	log.Println(state)
	log.Println("DEBUG Current:", state.Version, "Notification file:", notification.Version)
	if state.Version >= notification.Version {
		return
	}
	// -- compare with latest notification
	//    * is version newer? if not then blow
	//    * are there contiguous deltas since our last delta? if not, download snapshot
	//    * apply deltas
	ds.ApplyDeltas("RIR-TEST", []nrtm4model.Change{})
}

func validateNotificationFile(file nrtm4model.Notification) []error {
	var errs []error
	if file.NrtmVersion != 4 {
		errs = append(errs, newInvalidJSONError("notificationFile nrtm version is not v4: '%v'", file.NrtmVersion))
	}
	if len(file.SessionID) < 36 {
		errs = append(errs, newInvalidJSONError("notificationFile session ID is not valid: '%v'", file.SessionID))
	}
	if len(file.Source) < 3 {
		errs = append(errs, newInvalidJSONError("notificationFile source is not valid: '%v'", file.Source))
	}
	if file.Version < 1 {
		errs = append(errs, newInvalidJSONError("notificationFile version must be positive: '%v'", file.NrtmVersion))
	}
	return errs
}
