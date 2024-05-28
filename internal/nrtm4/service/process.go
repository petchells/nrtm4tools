package service

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/jsonseq"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/nrtm4model"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/rpsl"
)

var fileWriteBufferLength = 1024 * 8
var rpslInsertBatchSize = 1000

// NewNRTMProcessor injects repo and client into service and return a new instance
func NewNRTMProcessor(config AppConfig, repo persist.Repository, client Client) NRTMProcessor {
	return NRTMProcessor{
		config: config,
		repo:   repo,
		client: client,
	}
}

// NRTMProcessor orchestration for functions the client implements
type NRTMProcessor struct {
	config AppConfig
	repo   persist.Repository
	client Client
}

// Connect stores details about a connection
func (p NRTMProcessor) Connect(notificationURL string, label string) error {
	dl := downloader{}
	notification, errs := dl.downloadNotificationFile(p.client, notificationURL)
	if len(errs) > 0 {
		return errors.New("download error(s): " + errs[0].Error())
	}
	ds := NrtmDataService{Repository: p.repo}
	if ds.getSourceByURLAndLabel(notificationURL, label) != nil {
		return errors.New("source already exists")
	}
	// Download snapshot
	snapshotFile, err := p.storeSnapshotInEmptyDirectory(notification.SnapshotRef.URL)
	if err != nil {
		return err
	}
	defer snapshotFile.Close()

	source := persist.NewNRTMSource(notification, "")
	if source, err = ds.saveNewSource(source); err != nil {
		logger.Error("saving new source. Remove Source and restart sync", err)
		return err
	}
	if err = p.insertSnapshotRecords(source, snapshotFile); err != nil {
		logger.Error("when inserting snapshot records. Remove Source and restart sync", err)
		return err
	}
	return nil
}

func (p NRTMProcessor) storeSnapshotInEmptyDirectory(snapshotURL string) (*os.File, error) {
	fm := fileManager{client: p.client}
	snapshotOSFile, err := os.Open(filepath.Join(p.config.NRTMFilePath, filepath.Base(snapshotURL)))
	if err != nil {
		if err = os.RemoveAll(p.config.NRTMFilePath); err != nil {
			logger.Warn("removed existing directory", err)
		}
		err = os.Mkdir(p.config.NRTMFilePath, 0755)
		if err != nil {
			log.Fatal(err)
		}
		logger.Info("Created path", "path", p.config.NRTMFilePath)
		logger.Info("Downloading snapshot", "url", snapshotURL)
		if snapshotOSFile, err = fm.writeResourceToPath(snapshotURL, p.config.NRTMFilePath); err != nil {
			log.Fatal(err)
		}
		if err != nil {
			logger.Error("failed to write snapshot", "url", snapshotURL, "path", p.config.NRTMFilePath)
			return nil, err
		}
	}
	return snapshotOSFile, err
}

func (p NRTMProcessor) insertSnapshotRecords(source persist.NRTMSource, snapshotOSFile *os.File) error {
	defer snapshotOSFile.Close()
	fm := fileManager{client: p.client}
	if err := fm.readSnapshotRecords(snapshotOSFile, snapshotObjectInsertionFunc(p.repo, source)); err != io.EOF {
		logger.Warn("Failed to initialize source", "source", source, err)
		return err
	}
	return nil
}

// ListSources shows all sources
func (p NRTMProcessor) ListSources() ([]persist.NRTMSource, error) {
	ds := NrtmDataService{Repository: p.repo}
	return ds.getSources()
}

// UpdateNRTM updates the repo using data fetched from the client at the given url, storing files in nrtmFilePath
// func (p NRTMProcessor) UpdateNRTM(source string, label string) {
// 	ds := NrtmDataService{Repository: p.repo}
// 	src, err := ds.fetchSource(source, label)
// 	if err != nil {
// 		log.Println("ERROR UpdateNRTM", err)
// 		return
// 	}

// 	dl := downloader{}
// 	notification, err := dl.downloadNotificationFile(p.client, src.URL)
// 	if err != nil {
// 		log.Println("ERROR UpdateNRTM", err)
// 		return
// 	}
// 	state, clientErr := repo.GetState(notification.Source)
// 	if clientErr == &persist.ErrStateNotInitialized {
// 		log.Println("INFO No previous state found. Initializing")
// 		fileName, err := fileNameFromURLString(url)
// 		if err != nil {
// 			log.Println("ERROR URL:", url, err)
// 			return
// 		}
// 		state := persist.NRTMFile{
// 			ID:           0,
// 			Created:      time.Now().UTC(),
// 			NrtmSourceID: notification.Source,
// 			Version:      notification.Version,
// 			URL:          url,
// 			Type:         persist.NotificationFile,
// 			FileName:     fileName,
// 		}
// 		if err = repo.SaveFile(&state); err != nil {
// 			log.Println("ERROR Saving state", err)
// 			return
// 		}
// 		fm := fileManager{repo: repo, client: client}
// 		// save notification file to disk and nrtmstate table
// 		snapshotOSFile, err := os.Open(filepath.Join(nrtmFilePath, filepath.Base(notification.Snapshot.Url)))
// 		if err != nil {
// 			if err = os.RemoveAll(nrtmFilePath); err != nil {
// 				log.Println("WARNING removed existing directory", err)
// 			}
// 			err = os.Mkdir(nrtmFilePath, 0755)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			log.Println("INFO Created path", nrtmFilePath)
// 			log.Println("INFO Downloading snapshot", notification.Snapshot.Url)
// 			if snapshotOSFile, err = fm.writeResourceToPath(notification.Snapshot.Url, nrtmFilePath); err != nil {
// 				log.Fatal(err)
// 			}
// 			if err != nil {
// 				log.Println("ERROR failed to write snapshot", notification.Snapshot.Url, "to", nrtmFilePath)
// 				return
// 			}
// 		}
// 		defer snapshotOSFile.Close()

// 		if err = fm.readSnapshotRecords(snapshotOSFile, snapshotRecordReaderFunc(repo, state)); err != io.EOF {
// 			log.Println("WARN failed to initialize source", state, err)
// 			return
// 		}
// 		var stateErr error
// 		if state, stateErr = repo.GetState(notification.Source); stateErr != nil {
// 			log.Println("ERROR failed to retrieve initial state", stateErr)
// 			return
// 		}
// 		log.Println("INFO new state", state)
// 	} else if clientErr != nil {
// 		log.Println("ERROR Database error", clientErr)
// 		return
// 	}

// 	// -- compare with latest notification
// 	//    * is version newer? if not then blow
// 	//    * are there contiguous deltas since our last delta? if not, download snapshot
// 	//    * apply deltas
// 	log.Println("DEBUG Current:", state.Version, "Notification file:", notification.Version)
// 	if state.Version >= notification.Version {
// 		log.Println("INFO Nothing to do: version is up to date.")
// 		return
// 	}
// 	log.Println("DEBUG Applying deltas >", state.Version, "up to", notification.Version)
// 	// TODO:
// 	// Get the actual deltas
// 	ds.applyDeltas(notification.Source, []nrtm4model.Change{})
// }

func snapshotObjectInsertionFunc(repo persist.Repository, source persist.NRTMSource) jsonseq.RecordReaderFunc {

	successfulObjects := 0
	failedObjects := 0

	var rpslObjects []rpsl.Rpsl

	return func(bytes []byte, err error) error {
		if err == &persist.ErrNoEntity {
			logger.Warn("error empty JSON", err)
			return err
		}
		if err == io.EOF {
			// Expected error: end of snapshot objects. Round them up and save them.
			so := new(nrtm4model.SnapshotObjectJSON)
			if err = json.Unmarshal(bytes, so); err == nil {
				rpsl, err := rpsl.ParseString(so.Object)
				if err != nil {
					return err
				}
				successfulObjects++
				rpslObjects = append(rpslObjects, rpsl)
				return repo.SaveSnapshotObjects(source, rpslObjects)
			}
			return err
		} else if err != nil {
			// Unexpected error. Should be able to read snapshot header or objects.
			logger.Warn("error unmarshalling JSON.", err)
			return err
		} else if successfulObjects == 0 {
			// First record is the Snapshot header
			successfulObjects++
			sf := new(nrtm4model.SnapshotFileJSON)
			if err = json.Unmarshal(bytes, sf); err != nil {
				logger.Warn("error unmarshalling JSON. Expected SnapshotFile", err, "errors", failedObjects)
				return err
			}
			return nil
		} else {
			// Subsequent records are objects
			so := new(nrtm4model.SnapshotObjectJSON)
			if err = json.Unmarshal(bytes, so); err == nil {
				rpsl, err := rpsl.ParseString(so.Object)
				if err != nil {
					failedObjects++
					return err
				}
				successfulObjects++
				rpslObjects = append(rpslObjects, rpsl)
				if len(rpslObjects) >= rpslInsertBatchSize {
					err = repo.SaveSnapshotObjects(source, rpslObjects)
					if err != nil {
						return err
					}
					rpslObjects = nil
				}
				return nil
			}
			failedObjects++
			logger.Warn("Error unmarshalling JSON. Expected SnapshotObject", err, "numErrors", failedObjects)
			return err
		}
	}
}
