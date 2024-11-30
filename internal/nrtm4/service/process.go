package service

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"sort"
	"sync"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/jsonseq"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/rpsl"
)

var (
	// ErrNRTMVersionMismatch nrtm version is not 4
	ErrNRTMVersionMismatch = errors.New("nrtm version is not 4")
	// ErrNRTMSourceMismatch session id does not match source
	ErrNRTMSourceMismatch = errors.New("session id does not match source")
	// ErrNRTMSourceNameMismatch source name does not match source
	ErrNRTMSourceNameMismatch = errors.New("source name does not match source")
	// ErrNRTMFileVersionMismatch file version does not match its reference
	ErrNRTMFileVersionMismatch = errors.New("file version does not match its reference")
	// ErrNRTMFileVersionInconsistency version is lower than source
	ErrNRTMFileVersionInconsistency = errors.New("version is lower than source")

	fileWriteBufferLength = 1024 * 8
	rpslInsertBatchSize   = 1000
)

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
	if source, err = ds.saveNewSource(source, notification); err != nil {
		logger.Error("There was a problem saving the source. Remove it and restart sync", "error", err)
		return err
	}
	logger.Info("Inserting snapshot objects")
	if err := fm.readJSONSeqRecords(snapshotFile, snapshotObjectInsertFunc(p.repo, source, notification)); err != io.EOF {
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
		if err := fm.readJSONSeqRecords(file, applyDeltaFunc(p.repo, source, notification, deltaRef)); err != io.EOF {
			logger.Warn("Failed to apply delta", "source", source, "error", err)
			return err
		}
	}
	logger.Info("Finished syncing deltas")
	return nil
}

func applyDeltaFunc(repo persist.Repository, source persist.NRTMSource, notification persist.NotificationJSON, deltaRef persist.FileRefJSON) jsonseq.RecordReaderFunc {
	var header *persist.DeltaFileJSON
	return func(bytes []byte, err error) error {
		if err == &persist.ErrNoEntity {
			logger.Warn("error empty JSON", "error", err)
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
				_, err = repo.SaveSource(source, notification)
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
				return errors.New("no delta action available: " + delta.Action)
			}
			return nil
		}
		return err
	}
}

// Counter is a counter with a mutex
type Counter struct {
	mu sync.Mutex
	n  int64
}

// Increment increments the counter
func (c *Counter) Increment() {
	c.mu.Lock()
	c.n++
	c.mu.Unlock()
}

// RPSLObjectList an ummutable list of objects
type RPSLObjectList struct {
	//mu      sync.Mutex
	objects []rpsl.Rpsl
}

// NewRPSLObjectList returns an initialized RPSLObjectList
func NewRPSLObjectList() RPSLObjectList {
	return RPSLObjectList{make([]rpsl.Rpsl, rpslInsertBatchSize, rpslInsertBatchSize*2)}
}

// Add adds an object the list
func (l *RPSLObjectList) Add(obj rpsl.Rpsl) {
	//l.mu.Lock()
	l.objects = append(l.objects, obj)
	//l.mu.Unlock()
}

// GetBatch will return a slice of objects only if 'size' are available. They are removed from the list
func (l *RPSLObjectList) GetBatch(size int) []rpsl.Rpsl {
	res := []rpsl.Rpsl{}
	//l.mu.Lock()
	if len(l.objects) >= size {
		res = l.objects[:size]
		l.objects = l.objects[size:]
	}
	//l.mu.Unlock()
	return res
}

// GetAll returns all RPSL objects and empties the internal list.
func (l *RPSLObjectList) GetAll() []rpsl.Rpsl {
	//l.mu.Lock()
	res := l.objects
	l.objects = []rpsl.Rpsl{}
	//l.mu.Unlock()
	return res
}

func snapshotObjectInsertFunc(repo persist.Repository, source persist.NRTMSource, notification persist.NotificationJSON) jsonseq.RecordReaderFunc {

	//	var rpslObjects []rpsl.Rpsl
	var snapshotHeader *persist.SnapshotFileJSON

	var wg sync.WaitGroup

	//objectCh := make(chan rpsl.Rpsl, 2000)
	objectList := RPSLObjectList{}
	//objectList := make([]rpsl.Rpsl, rpslInsertBatchSize, rpslInsertBatchSize*2)
	successfulObjects := Counter{}
	failedObjects := Counter{}

	// unmarshalBytesToChan := func(bytes []byte) {
	// 	so := new(persist.SnapshotObjectJSON)
	// 	if err := json.Unmarshal(bytes, so); err == nil {
	// 		rpsl, err := rpsl.ParseString(so.Object)
	// 		if err != nil {
	// 			failedObjects.Increment()
	// 			logger.Warn("Failed to parse rpsl.Rpsl from", "so.Object", so.Object, "error", err)
	// 		}
	// 		objectCh <- rpsl
	// 		successfulObjects.Increment()
	// 	} else {
	// 		logger.Warn("Failed to unmarshal RPSL string from", "so.Object", so.Object, "error", err)
	// 	}
	// }

	// saveObjects := func(ch chan rpsl.Rpsl) {
	// 	for {
	// 		select {
	// 		case rpsl := <-ch:
	// 			objectListSync.Add(rpsl)
	// 			rpslObjects := objectListSync.GetBatch(rpslInsertBatchSize)
	// 			if len(rpslObjects) > 0 {
	// 				wg.Add(1)
	// 				go func() {
	// 					defer wg.Done()
	// 					err := repo.SaveSnapshotObjects(source, rpslObjects, snapshotHeader.NrtmFileJSON)
	// 					if err != nil {
	// 						log.Fatalln("Error saving snapshot object", err)
	// 					}
	// 				}()
	// 			}
	// 		default:
	// 		}
	// 	}
	// }

	return func(bytes []byte, err error) error {
		if err == &persist.ErrNoEntity {
			logger.Warn("error empty JSON", "error", err)
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
				successfulObjects.Increment()
				objectList.Add(rpsl)
				rpslObjects := objectList.GetAll()
				wg.Wait()
				err = repo.SaveSnapshotObjects(source, rpslObjects, snapshotHeader.NrtmFileJSON)
				if err != nil {
					return err
				}
				source.Version = snapshotHeader.Version
				_, err = repo.SaveSource(source, notification)
				return err
			}
			return err
		} else if err != nil {
			// Unexpected error. Should be able to read snapshot header or objects.
			logger.Warn("error unmarshalling JSON.", "error", err)
			return err
		} else if successfulObjects.n == 0 {
			// First record is the Snapshot header
			successfulObjects.Increment()
			sf := new(persist.SnapshotFileJSON)
			if err = json.Unmarshal(bytes, sf); err != nil {
				logger.Warn("error unmarshalling JSON. Expected SnapshotFile", "error", err, "numFailures", failedObjects.n)
				return err
			}
			if sf.Version != notification.SnapshotRef.Version {
				return ErrNRTMFileVersionMismatch
			}
			snapshotHeader = sf
			return nil
		} else {
			// Subsequent records are objects
			so := new(persist.SnapshotObjectJSON)
			if err := json.Unmarshal(bytes, so); err == nil {
				rpsl, err := rpsl.ParseString(so.Object)
				if err != nil {
					failedObjects.Increment()
					logger.Warn("Failed to parse rpsl.Rpsl from", "so.Object", so.Object, "error", err)
				}
				objectList.Add(rpsl)
				successfulObjects.Increment()
			} else {
				logger.Warn("Failed to unmarshal RPSL string from", "so.Object", so.Object, "error", err)
			}
			rpslObjects := objectList.GetBatch(rpslInsertBatchSize)
			if len(rpslObjects) > 0 {
				err := repo.SaveSnapshotObjects(source, rpslObjects, snapshotHeader.NrtmFileJSON)
				if err != nil {
					log.Fatalln("Error saving snapshot object", err)
				}
			}
			return nil
		}
	}

}

func validateDeltaHeader(file persist.NrtmFileJSON, source persist.NRTMSource, deltaRef persist.FileRefJSON) error {
	if file.NrtmVersion != 4 {
		return ErrNRTMVersionMismatch
	}
	if file.SessionID != source.SessionID {
		return ErrNRTMSourceMismatch
	}
	if file.Source != source.Source {
		return ErrNRTMSourceNameMismatch
	}
	if file.Version != deltaRef.Version {
		return ErrNRTMFileVersionMismatch
	}
	if file.Version < source.Version {
		return ErrNRTMFileVersionInconsistency
	}
	return nil
}
