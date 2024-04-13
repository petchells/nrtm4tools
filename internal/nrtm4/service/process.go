package service

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/jsonseq"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/nrtm4model"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/rpsl"
)

var fileBufferLength = 1024 * 8

// UpdateNRTM updates the repo using data fetched from the client at the given url, storing files in nrtmFilePath
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
	state, clientErr := repo.GetState(notification.Source)
	if clientErr == &persist.ErrStateNotInitialized {
		log.Println("INFO No previous state found. Initializing")
		fileName, err := fileNameFromURLString(url)
		if err != nil {
			log.Println("ERROR URL:", url, err)
			return
		}
		state := persist.NRTMState{
			ID:       0,
			Created:  time.Now(),
			Source:   notification.Source,
			Version:  notification.Version,
			URL:      url,
			Type:     persist.NotificationFile,
			FileName: fileName,
		}
		if err = repo.SaveState(&state); err != nil {
			log.Println("ERROR Saving state", err)
			return
		}
		fm := fileManager{repo: repo, client: client}
		// save notification file to disk and nrtmstate table
		snapshotOSFile, err := os.Open(filepath.Join(nrtmFilePath, filepath.Base(notification.Snapshot.Url)))
		if err != nil {
			if err = os.RemoveAll(nrtmFilePath); err != nil {
				log.Println("WARNING removed existing directory", err)
			}
			err = os.Mkdir(nrtmFilePath, 0755)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("INFO Created path", nrtmFilePath)
			log.Println("INFO Downloading snapshot", notification.Snapshot.Url)
			if snapshotOSFile, err = fm.writeResourceToPath(notification.Snapshot.Url, nrtmFilePath); err != nil {
				log.Fatal(err)
			}
			if err != nil {
				log.Println("ERROR failed to write snapshot", notification.Snapshot.Url, "to", nrtmFilePath)
				return
			}
		}
		defer snapshotOSFile.Close()

		if err = fm.readSnapshotRecords(snapshotOSFile, snapshotRecordReaderFunc(repo, state)); err != io.EOF {
			log.Println("WARN failed to initialize source", state, err)
			return
		}
		var stateErr error
		if state, stateErr = repo.GetState(notification.Source); stateErr != nil {
			log.Println("ERROR failed to retrieve initial state", stateErr)
			return
		}
		log.Println("INFO new state", state)
	} else if clientErr != nil {
		log.Println("ERROR Database error", clientErr)
		return
	}

	// -- compare with latest notification
	//    * is version newer? if not then blow
	//    * are there contiguous deltas since our last delta? if not, download snapshot
	//    * apply deltas
	log.Println("DEBUG Current:", state.Version, "Notification file:", notification.Version)
	if state.Version >= notification.Version {
		log.Println("INFO Nothing to do: version is up to date.")
		return
	}
	log.Println("DEBUG Applying deltas >", state.Version, "up to", notification.Version)
	// TODO:
	// Get the actual deltas
	ds.applyDeltas(notification.Source, []nrtm4model.Change{})
}

func snapshotRecordReaderFunc(repo persist.Repository, state persist.NRTMState) jsonseq.RecordReaderFunc {

	i := 0
	failedEntities := 0

	var rpslObjects []rpsl.Rpsl

	return func(bytes []byte, err error) error {
		if err != &persist.ErrNoEntity {
			if err == io.EOF {
				so := new(nrtm4model.SnapshotObject)
				if err = json.Unmarshal(bytes, so); err == nil {
					rpsl, err := rpsl.ParseString(so.Object)
					if err != nil {
						return err
					}
					i++
					rpslObjects = append(rpslObjects, rpsl)
					return repo.SaveSnapshotObjects(state, rpslObjects)
				}
				return err
			} else if err != nil {
				log.Println("WARN error unmarshalling JSON.", err)
				return err
			} else if i == 0 {
				sf := new(nrtm4model.SnapshotFile)
				if err = json.Unmarshal(bytes, sf); err == nil {
					return repo.SaveSnapshotFile(state, *sf)
				} else {
					log.Println("WARN error unmarshalling JSON. Expected SnapshotFile", err, "errors", failedEntities)
					return err
				}
			} else {
				so := new(nrtm4model.SnapshotObject)
				if err = json.Unmarshal(bytes, so); err == nil {
					rpsl, err := rpsl.ParseString(so.Object)
					if err != nil {
						return err
					}
					i++
					rpslObjects = append(rpslObjects, rpsl)
					if len(rpslObjects) >= 1000 {
						err = repo.SaveSnapshotObjects(state, rpslObjects)
						if err != nil {
							return err
						}
						rpslObjects = nil
					}
					return nil
				}
				failedEntities++
				log.Println("WARN error unmarshalling JSON. Expected SnapshotObject", err, "errors", failedEntities)
				return err

			}
		} else {
			log.Println("WARN error empty JSON", err)
			return err
		}
	}
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
