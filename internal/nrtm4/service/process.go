package service

import (
	"io"
	"log"
	"os"
	"strings"
	"time"

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
	ds := NrtmDataService{Repository: repo}
	state, err := repo.GetState(notification.Source)
	if err != nil {
		log.Println("Failed to get state", err)
		//repo.CreateState(state)
		// -- if no state, then initialize
		//    * get snapshot, put file on disk
		//    * parse it
		//    * save state
		//    * insert rpsl objects
		//    * see if there are more deltas to process
		//    * done and dusted
		var snapshotFile *os.File
		if snapshotFile, err = writeResourceToPath(notification.Snapshot.Url, nrtmFilePath); err != nil {
			log.Println("ERROR occurred when writing snapshot file to disk", notification.Snapshot.Url, err)
			return
		}
		defer func() {
			if err := snapshotFile.Close(); err != nil {
				panic(err)
			}
		}()
		state = persist.NRTMState{
			ID:      0,
			Created: time.Now(),
			Source:  notification.Source,
			Version: notification.Version,
			URL:     url,
			Type:    persist.Notification,
			Payload: "",
		}
		err = repo.SaveState(state)
		if err != nil {
			log.Println("WARN failed to save state", state)
			return
		}
		log.Println("DEBUG snapshotFile.Name()", snapshotFile.Name())

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

func writeResourceToPath(url string, path string) (*os.File, error) {
	fileName := url[strings.LastIndex(url, "/")+1:]
	var reader io.ReadCloser
	var httpClient HttpClient
	var outFile *os.File
	var err error
	if reader, err = httpClient.fetchFile(url); err != nil {
		log.Println("ERROR Failed to fetch file", url, err)
		return nil, err
	}
	if outFile, err = os.Create(path + "/" + fileName); err != nil {
		log.Println("ERROR Failed to open file on disk", err)
		return nil, err
	}
	if err = transferReaderToFile(reader, outFile); err != nil {
		log.Println("ERROR writing file")
		return nil, err
	}
	return outFile, err
}

func transferReaderToFile(from io.ReadCloser, to *os.File) error {
	buf := make([]byte, FILE_BUFFER_LENGTH)
	for {
		n, err := from.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		if _, err := to.Write(buf[:n]); err != nil {
			return err
		}
	}
	return nil
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
