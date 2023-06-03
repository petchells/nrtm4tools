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
	httpClient := newNrtmHttpClient(url)
	notificationFile, err := httpClient.getUpdateNotification()
	if err != nil {
		log.Println("Failed to get notification file", err)
		return
	}
	if errs := validateNotificationFile(notificationFile); len(errs) > 0 {
		for _, err := range errs {
			log.Println("ERROR failed to get notificationFile", err)
		}
		return
	}
	// Fetch state
	ds := NrtmDataService{Repository: repo}
	state, err := repo.GetState(notificationFile.Source)
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

func validateNotificationFile(notificationFile nrtm4model.Notification) []error {
	var errs []error
	if notificationFile.NrtmVersion != 4 {
		errs = append(errs, newInvalidJSONError("notificationFile nrtm version is not v4: '%v'", notificationFile.NrtmVersion))
	}
	if len(notificationFile.SessionID) < 36 {
		errs = append(errs, newInvalidJSONError("notificationFile session ID is not valid: '%v'", notificationFile.SessionID))
	}
	if len(notificationFile.Source) < 3 {
		errs = append(errs, newInvalidJSONError("notificationFile source is not valid: '%v'", notificationFile.Source))
	}
	if notificationFile.Version < 1 {
		errs = append(errs, newInvalidJSONError("notificationFile version must be positive: '%v'", notificationFile.NrtmVersion))
	}
	return errs
}
