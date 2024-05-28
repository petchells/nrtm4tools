package service

import (
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/nrtm4model"
)

type downloader struct{}

func (dl downloader) downloadNotificationFile(client Client, url string) (nrtm4model.NotificationJSON, []error) {
	var notification nrtm4model.NotificationJSON
	var err error
	if notification, err = client.getUpdateNotification(url); err != nil {
		logger.Error("fetching notificationFile", err)
		return notification, []error{err}
	}
	errs := dl.validateNotificationFile(notification)
	return notification, errs
}

func (dl downloader) validateNotificationFile(file nrtm4model.NotificationJSON) []error {
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
	if len(file.SnapshotRef.URL) < 20 {
		errs = append(errs, newNRTMServiceError("notificationFile snapshot url is not valid: '%v'", file.SnapshotRef.URL))
	}
	return errs
}
