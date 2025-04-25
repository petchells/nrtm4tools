package service

import (
	"encoding/json"
	"errors"
	"io"
	"path/filepath"
	"sort"

	"github.com/jackc/pgx/v5"
	"github.com/petchells/nrtm4tools/internal/nrtm4/jsonseq"
	"github.com/petchells/nrtm4tools/internal/nrtm4/persist"
	"github.com/petchells/nrtm4tools/internal/nrtm4/rpsl"
)

func syncDeltas(p NRTMProcessor, notification persist.NotificationJSON, source persist.NRTMSource) (persist.NRTMSource, error) {
	dlDir := filepath.Join(p.config.NRTMFilePath, source.Source, source.SessionID)
	deltaRefs, err := findUpdates(notification, source)
	if err != nil {
		return source, err
	}
	fm := fileManager{p.client}
	ds := NrtmDataService{Repository: p.repo}
	for _, deltaRef := range deltaRefs {
		UserLogger.Info("Fetching delta", "version", deltaRef.Version, "relurl", deltaRef.URL)
		file, err := fm.fetchFileAndCheckHash(source.NotificationURL, deltaRef, dlDir)
		if err != nil {
			UserLogger.Error("Error fetching delta", "source", source.Source, "delta", deltaRef.Version, "relurl", deltaRef.URL, "error", err)
			return source, err
		}
		defer file.Close()
		if err := fm.readJSONSeqRecords(file, applyDeltaFunc(p.repo, source, deltaRef)); err != io.EOF {
			UserLogger.Error("Failed to apply delta", "source", source.Source, "delta", deltaRef.Version, "relurl", deltaRef.URL)
			return source, err
		}
		source.Version = uint32(deltaRef.Version)
		src, err := ds.saveSource(source)
		if err != nil {
			return source, err
		}
		source = *src
	}
	UserLogger.Info("Delta sync complete", "number of deltas files applied", len(deltaRefs))
	return source, nil
}

func findUpdates(notification persist.NotificationJSON, source persist.NRTMSource) ([]persist.FileRefJSON, error) {

	deltaRefs := []persist.FileRefJSON{}
	for _, deltaRef := range notification.DeltaRefs {
		if deltaRef.Version > int64(source.Version) {
			deltaRefs = append(deltaRefs, deltaRef)
		}
	}
	if len(deltaRefs) == 0 {
		return nil, nil
	}
	sort.Slice(deltaRefs, func(r1, r2 int) bool {
		return deltaRefs[r1].Version < deltaRefs[r2].Version
	})
	if source.Version+1 < uint32(deltaRefs[0].Version) {
		return nil, ErrNextConsecutiveDeltaUnavaliable
	}
	return deltaRefs, nil
}

func applyDeltaFunc(repo persist.Repository, source persist.NRTMSource, deltaRef persist.FileRefJSON) jsonseq.RecordReaderFunc {
	var header *persist.DeltaFileJSON
	return func(bytes []byte, err error) error {
		if err != nil && err != io.EOF { // eof also gives us a record
			return err
		}
		if header == nil {
			deltaHeader := new(persist.DeltaFileJSON)
			if err = json.Unmarshal(bytes, deltaHeader); err != nil {
				return err
			}
			if err = validateDeltaHeader(deltaHeader.NrtmFileJSON, source, deltaRef); err != nil {
				return err
			}
			header = deltaHeader
			return err
		}
		delta := new(persist.DeltaJSON)
		if err = json.Unmarshal(bytes, delta); err != nil {
			return err
		}
		switch {
		case delta.Action == persist.DeltaAddModifyAction:
			rpsl, err := rpsl.ParseFromJSONString(*delta.Object)
			if err != nil {
				UserLogger.Error("Cannot parse RPSL for AddModify action", "object", *delta.Object, "error", err)
				return err
			}
			err = repo.AddModifyObject(source, rpsl, header.NrtmFileJSON)
			if err != nil {
				UserLogger.Error("Delta AddModifyObject failed", "rpsl", rpsl, "relurl", deltaRef.URL, "error", err)
				return err
			}
		case delta.Action == persist.DeltaDeleteAction:
			err = repo.DeleteObject(source, *delta.ObjectClass, *delta.PrimaryKey, header.NrtmFileJSON)
			if err != nil {
				if err == pgx.ErrNoRows {
					const txt = "Delta delete_object failed because object is not in the repository"
					UserLogger.Error(txt, "url", deltaRef.URL, "ObjectClass", *delta.ObjectClass, "PrimaryKey", *delta.PrimaryKey)
					return newNRTMServiceError("%v. class: %v primary-key: %v", txt, *delta.ObjectClass, *delta.PrimaryKey)
				}
				UserLogger.Error("Delta DeleteObject failed", "url", deltaRef.URL, "ObjectClass", *delta.ObjectClass, "PrimaryKey", *delta.PrimaryKey, "error", err)
				return err
			}

		default:
			UserLogger.Error("Delta file contains invalid action", "url", deltaRef.URL, "delta.Action", delta.Action)
			return errors.New("invalid delta action")
		}
		return nil
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
	if file.Version < int64(source.Version) {
		return ErrNRTM4FileVersionInconsistency
	}
	return nil
}
