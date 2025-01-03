package service

import (
	"fmt"
	"log"
	"strings"

	"github.com/petchells/nrtm4client/internal/nrtm4/persist"
)

// ErrNRTMServiceError is when sth is wrong with the NRTM server
type ErrNRTMServiceError struct {
	Message string
}

func (e ErrNRTMServiceError) Error() string {
	return "nrtm service error: " + e.Message
}

func newNRTMServiceError(msg string, args ...any) ErrNRTMServiceError {
	return ErrNRTMServiceError{fmt.Sprintf(msg, args...)}
}

// NrtmDataService is an implementation of a persist.Repository
type NrtmDataService struct {
	Repository persist.Repository
}

func (ds NrtmDataService) getSourceByURLAndLabel(url string, label string) *persist.NRTMSource {
	sources, err := ds.getSources()
	if err != nil {
		log.Panicln("Failure calling NrtmDataService.getSourceByURLAndLabel", err)
	}
	for _, src := range sources {
		if strings.EqualFold(src.NotificationURL, url) && strings.EqualFold(src.Label, label) {
			found := src
			return &found
		}
	}
	return nil
}

func (ds NrtmDataService) getSourceByNameAndLabel(name string, label string) *persist.NRTMSource {
	sources, err := ds.getSources()
	if err != nil {
		log.Panicln("Failure calling NrtmDataService.getSourceByNameAndLabel", err)
	}
	for _, src := range sources {
		if strings.EqualFold(src.Source, name) && strings.EqualFold(src.Label, label) {
			found := src
			return &found
		}
	}
	return nil
}

func (ds NrtmDataService) getSources() ([]persist.NRTMSource, error) {
	return ds.Repository.GetSources()
}

func (ds NrtmDataService) getNotifications(src persist.NRTMSource, from, to uint32) ([]persist.Notification, error) {
	return ds.Repository.GetNotificationHistory(src, from, to)
}

func (ds NrtmDataService) saveNewSource(source persist.NRTMSource, notification persist.NotificationJSON) (persist.NRTMSource, error) {
	return ds.Repository.SaveSource(source, notification)
}
