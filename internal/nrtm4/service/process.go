package service

import (
	"encoding/json"
	"errors"
	"io"
	"sort"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/jsonseq"
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
	fm := fileManager{p.client}
	logger.Info("Fetching notification")
	notification, errs := fm.downloadNotificationFile(notificationURL)
	if len(errs) > 0 {
		return errors.New("download error(s): " + errs[0].Error())
	}
	ds := NrtmDataService{Repository: p.repo}
	if ds.getSourceByURLAndLabel(notificationURL, label) != nil {
		return errors.New("source already exists")
	}
	err := fm.ensureDirectoryExists(p.config.NRTMFilePath)
	if err != nil {
		return err
	}
	// Download snapshot
	logger.Info("Fetching snapshot file...")
	snapshotFile, err := fm.fetchFileAndCheckHash(notification.SnapshotRef, p.config.NRTMFilePath)
	if err != nil {
		return err
	}
	logger.Info("Snapshot file downloaded")
	defer snapshotFile.Close()

	logger.Info("Saving new source", "source", notification.Source)
	source := persist.NewNRTMSource(notification, label, notificationURL)
	if source, err = ds.saveNewSource(source); err != nil {
		logger.Error("There was a problem saving the source. Remove it and restart sync", "error", err)
		return err
	}
	logger.Info("Inserting snapshot objects")
	if err := fm.readJSONSeqRecords(snapshotFile, snapshotObjectInsertFunc(p.repo, source, notification.SnapshotRef)); err != io.EOF {
		logger.Error("Invalid snapshot. Remove Source and restart sync", "error", err)
		return err
	}
	return p.syncDeltas(notification, source)
}

// Update brings the local mirror up to date
func (p NRTMProcessor) Update(sourceName string, label string) error {
	ds := NrtmDataService{Repository: p.repo}
	source := ds.getSourceByNameAndLabel(sourceName, label)
	if source == nil {
		logger.Warn("No source with given name and label", "name", sourceName, "label", label)
		return errors.New("no source found")
	}
	fm := fileManager{p.client}
	notification, errs := fm.downloadNotificationFile(source.NotificationURL)
	if len(errs) > 0 {
		for _, e := range errs {
			logger.Error("Problem downloading notification file", "error", e)
		}
		return errors.New("problem downloading notification file")
	}
	if notification.SessionID != source.SessionID {
		return errors.New("server has a new mirror session")
	}
	if notification.Version < source.Version {
		return errors.New("server has old version")
	}
	if notification.Version == source.Version {
		logger.Info("Already at latest version")
		return nil
	}
	return p.syncDeltas(notification, *source)
}

// ListSources shows all sources
func (p NRTMProcessor) ListSources() ([]persist.NRTMSource, error) {
	ds := NrtmDataService{Repository: p.repo}
	return ds.getSources()
}

func (p NRTMProcessor) syncDeltas(notification persist.NotificationJSON, source persist.NRTMSource) error {
	logger.Info("Looking for deltas")
	deltaRefs := []persist.FileRefJSON{}
	for _, deltaRef := range *notification.DeltaRefs {
		if deltaRef.Version > source.Version {
			deltaRefs = append(deltaRefs, deltaRef)
		}
	}
	if len(deltaRefs) == 0 {
		return errors.New("restart sync mirror is too old")
	}
	logger.Info("Found deltas", "numdeltas", len(deltaRefs))
	sort.Sort(fileRefsByVersion(deltaRefs))
	fm := fileManager{p.client}
	for _, deltaRef := range deltaRefs {
		logger.Info("Processing delta", "delta", deltaRef.Version, "url", deltaRef.URL)
		file, err := fm.fetchFileAndCheckHash(deltaRef, p.config.NRTMFilePath)
		if err != nil {
			return err
		}
		defer file.Close()
		if err := fm.readJSONSeqRecords(file, applyDeltaFunc(p.repo, source, deltaRef)); err != io.EOF {
			logger.Warn("Failed to apply delta", "source", source, err)
			return err
		}
	}
	logger.Info("Finished syncing deltas")
	return nil
}

func applyDeltaFunc(repo persist.Repository, source persist.NRTMSource, deltaRef persist.FileRefJSON) jsonseq.RecordReaderFunc {
	var header *persist.DeltaFileJSON
	return func(bytes []byte, err error) error {
		if err == &persist.ErrNoEntity {
			logger.Warn("error empty JSON", err)
			return err
		}
		if err == nil || err == io.EOF {
			if header == nil {
				deltaHeader := new(persist.DeltaFileJSON)
				if err = json.Unmarshal(bytes, deltaHeader); err != nil {
					return err
				}
				if err = validateDeltaHeader(deltaHeader.NrtmFileJSON, source, deltaRef); err != nil {
					return err
				}
				header = deltaHeader
				source.Version = deltaRef.Version
				_, err = repo.SaveSource(source)
				return err
			}
			delta := new(persist.DeltaJSON)
			if err = json.Unmarshal(bytes, delta); err != nil {
				return err
			}
			if delta.Action == persist.DeltaAddModifyAction {
				rpsl, err := rpsl.ParseString(*delta.Object)
				if err != nil {
					return err
				}
				repo.AddModifyObject(source, rpsl, header.NrtmFileJSON)
			} else if delta.Action == persist.DeltaDeleteAction {
				repo.DeleteObject(source, *delta.ObjectClass, *delta.PrimaryKey, header.NrtmFileJSON)
			} else {
				return errors.New("no action available: " + delta.Action)
			}
			return nil
		}
		return err
	}
}

func snapshotObjectInsertFunc(repo persist.Repository, source persist.NRTMSource, fileRef persist.FileRefJSON) jsonseq.RecordReaderFunc {

	successfulObjects := 0
	failedObjects := 0

	var rpslObjects []rpsl.Rpsl
	var snapshotHeader *persist.SnapshotFileJSON

	return func(bytes []byte, err error) error {
		if err == &persist.ErrNoEntity {
			logger.Warn("error empty JSON", err)
			return err
		}
		if err == io.EOF {
			// Expected error reading to end of snapshot objects. Round them up and save them.
			so := new(persist.SnapshotObjectJSON)
			if err = json.Unmarshal(bytes, so); err == nil {
				rpsl, err := rpsl.ParseString(so.Object)
				if err != nil {
					return err
				}
				successfulObjects++
				rpslObjects = append(rpslObjects, rpsl)
				err = repo.SaveSnapshotObjects(source, rpslObjects, snapshotHeader.NrtmFileJSON)
				if err != nil {
					return err
				}
				source.Version = snapshotHeader.Version
				_, err = repo.SaveSource(source)
				return err
			}
			return err
		} else if err != nil {
			// Unexpected error. Should be able to read snapshot header or objects.
			logger.Warn("error unmarshalling JSON.", err)
			return err
		} else if successfulObjects == 0 {
			// First record is the Snapshot header
			successfulObjects++
			sf := new(persist.SnapshotFileJSON)
			if err = json.Unmarshal(bytes, sf); err != nil {
				logger.Warn("error unmarshalling JSON. Expected SnapshotFile", err, "errors", failedObjects)
				return err
			}
			if sf.Version != fileRef.Version {
				return errors.New("snapshot header version does not match its reference")
			}
			snapshotHeader = sf
			return nil
		} else {
			// Subsequent records are objects
			so := new(persist.SnapshotObjectJSON)
			if err = json.Unmarshal(bytes, so); err == nil {
				rpsl, err := rpsl.ParseString(so.Object)
				if err != nil {
					failedObjects++
					return err
				}
				successfulObjects++
				rpslObjects = append(rpslObjects, rpsl)
				if len(rpslObjects) >= rpslInsertBatchSize {
					err = repo.SaveSnapshotObjects(source, rpslObjects, snapshotHeader.NrtmFileJSON)
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

func validateDeltaHeader(file persist.NrtmFileJSON, source persist.NRTMSource, deltaRef persist.FileRefJSON) error {
	if file.NrtmVersion != 4 {
		return errors.New("nrtm version is not 4")
	}
	if file.SessionID != source.SessionID {
		return errors.New("session id does not match source")
	}
	if file.Source != source.Source {
		return errors.New("source name does not match source")
	}
	if file.Version != deltaRef.Version {
		return errors.New("file version does not match its reference")
	}
	if file.Version < source.Version {
		return errors.New("version is lower than source")
	}
	return nil
}
