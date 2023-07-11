package service

import (
	"log"
	"os"

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
	state, err := repo.GetState(notification.Source)
	if err == &persist.ErrNoState {
		log.Println("INFO Failed to find previous state. Initializing")
		err = os.RemoveAll(nrtmFilePath)
		log.Println("WARNING removed existing directory", err)
		err = os.Mkdir(nrtmFilePath, 0755)
		if err != nil {
			log.Fatal(err)
		}

		fm := fileManager{repo: repo, client: client}
		// save notification file to disk and nrtmstate table
		var snapshotOSFile *os.File
		if snapshotOSFile, err = fm.writeResourceToPath(notification.Snapshot.Url, nrtmFilePath); err != nil {
			log.Fatal(err)
		}
		log.Println(snapshotOSFile.Name())

		i := 0
		if err = fm.initializeSourceAndParseSnapshot(url, nrtmFilePath, notification, func(bytes []byte, err error) {

			i++
		}); err != nil {
			log.Println("WARN failed to intialize source", state, err)
			return
		}
		if state, err = repo.GetState(notification.Source); err != nil {
			log.Println("ERROR failed to retrieve inital state", err)
			return
		}
		log.Println("INFO new state", state)
	} else if err != nil {
		log.Println("ERROR Database error", err)
		return
	}

	// -- compare with latest notification
	//    * is version newer? if not then blow
	//    * are there contiguous deltas since our last delta? if not, download snapshot
	//    * apply deltas
	log.Println("DEBUG Current:", state.Version, "Notification file:", notification.Version)
	if state.Version >= notification.Version {
		return
	}
	log.Println("DEBUG Applying deltas >", state.Version, "up to", notification.Version)
	ds := NrtmDataService{Repository: repo}
	// TODO:
	// Get the actual deltas
	ds.applyDeltas(notification.Source, []nrtm4model.Change{})
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
