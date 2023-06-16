package service

import (
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/nrtm4model"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
)

var FILE_BUFFER_LENGTH = 1024 * 8

func UpdateNRTM(repo persist.Repository, client Client, url string, nrtmFilePath string) {
	// Fetch notification
	// -- validate
	// -- new version?
	var notification nrtm4model.Notification
	var err error

	if notification, err = client.getUpdateNotification(url); err != nil {
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
	//repo.CreateState(state)
	// -- if no state, then initialize
	//    * get snapshot, put file on disk
	//    * parse it
	//    * save state
	//    * insert rpsl objects
	//    * see if there are more deltas to process
	//    * done and dusted
	ds := NrtmDataService{Repository: repo}
	state, err := repo.GetState(notification.Source)
	if err == pgx.ErrNoRows {
		log.Println("INFO Failed to find previous state. Initializing")
		var snapshotFile *os.File
		if snapshotFile, err = fileToDatabase(repo, notification.Snapshot.Url, notification.NrtmFile, persist.SnapshotFile, nrtmFilePath); err != nil {
			log.Println("WARN failed to save state", state, err)
			return
		}
		log.Println("DEBUG snapshotFile.Name()", snapshotFile.Name())
	} else if err != nil {
		log.Println("ERROR Database error", err)
		return
	}
	log.Println(state)
	// -- compare with latest notification
	//    * is version newer? if not then blow
	//    * are there contiguous deltas since our last delta? if not, download snapshot
	//    * apply deltas
	log.Println("DEBUG Current:", state.Version, "Notification file:", notification.Version)
	if state.Version >= notification.Version {
		return
	}
	ds.ApplyDeltas(notification.Source, []nrtm4model.Change{})
}

func validateNotificationFile(file nrtm4model.Notification) []error {
	var errs []error
	if file.NrtmVersion != 4 {
		errs = append(errs, newNRTMServiceError("notificationFile nrtm version is not v4: '%v'", file.NrtmVersion))
	}
	if len(file.SessionID) < 36 {
		errs = append(errs, newNRTMServiceError("notificationFile session ID is not valid: '%v'", file.SessionID))
	}
	if len(file.Source) < 3 {
		errs = append(errs, newNRTMServiceError("notificationFile source is not valid: '%v'", file.Source))
	}
	if file.Version < 1 {
		errs = append(errs, newNRTMServiceError("notificationFile version must be positive: '%v'", file.NrtmVersion))
	}
	if len(file.Snapshot.Url) < 20 {
		errs = append(errs, newNRTMServiceError("notificationFile snapshot url is not valid: '%v'", file.Snapshot.Url))
	}
	return errs
}
