package service

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"sort"

	"github.com/petchells/nrtm4client/internal/nrtm4/jsonseq"
	"github.com/petchells/nrtm4client/internal/nrtm4/persist"
	"github.com/petchells/nrtm4client/internal/nrtm4/rpsl"
)

func syncDeltas(p NRTMProcessor, notification persist.NotificationJSON, source persist.NRTMSource) error {
	deltaRefs, err := findUpdates(notification, source)
	if err != nil {
		return err
	}
	sort.Sort(fileRefsByVersion(deltaRefs))
	fm := fileManager{p.client}
	for _, deltaRef := range deltaRefs {
		logger.Info("Processing delta", "delta", deltaRef.Version, "url", deltaRef.URL)
		file, err := fm.fetchFileAndCheckHash(source.NotificationURL, deltaRef, p.config.NRTMFilePath)
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

func findUpdates(notification persist.NotificationJSON, source persist.NRTMSource) ([]persist.FileRefJSON, error) {

	if notification.DeltaRefs == nil || len(notification.DeltaRefs) == 0 {
		return nil, ErrNRTM4NoDeltasInNotification
	}

	deltaRefs := []persist.FileRefJSON{}
	versions := make([]uint32, len(notification.DeltaRefs))
	for i, deltaRef := range notification.DeltaRefs {
		versions[i] = deltaRef.Version
		if deltaRef.Version > source.Version {
			deltaRefs = append(deltaRefs, deltaRef)
		}
	}
	if len(deltaRefs) == 0 {
		log.Panic("Check the notification version before calling this function")
	}
	sort.Slice(deltaRefs, func(r1, r2 int) bool {
		return deltaRefs[r1].Version < deltaRefs[r2].Version
	})
	if source.Version+1 < deltaRefs[0].Version {
		return nil, ErrNextConsecutiveDeltaUnavaliable
	}
	logger.Info("Found deltas", "source", notification.Source, "numdeltas", len(deltaRefs))
	return deltaRefs, nil
}

func applyDeltaFunc(repo persist.Repository, source persist.NRTMSource, notification persist.NotificationJSON, deltaRef persist.FileRefJSON) jsonseq.RecordReaderFunc {
	var header *persist.DeltaFileJSON
	return func(bytes []byte, err error) error {
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
				rpsl, err := rpsl.ParseFromJSONString(*delta.Object)
				if err != nil {
					return err
				}
				err = repo.AddModifyObject(source, rpsl, header.NrtmFileJSON)
				if err != nil {
					logger.Error("Delta AddModifyO0bject failed", "rpsl", rpsl, "error", err)
					return err
				}
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

func validateDeltaHeader(file persist.NrtmFileJSON, source persist.NRTMSource, deltaRef persist.FileRefJSON) error {
	if file.NrtmVersion != 4 {
		return ErrNRTM4VersionMismatch
	}
	if file.SessionID != source.SessionID {
		return ErrNRTM4SourceMismatch
	}
	if file.Source != source.Source {
		return ErrNRTM4SourceNameMismatch
	}
	if file.Version != deltaRef.Version {
		return ErrNRTM4FileVersionMismatch
	}
	if file.Version < source.Version {
		return ErrNRTM4FileVersionInconsistency
	}
	return nil
}
