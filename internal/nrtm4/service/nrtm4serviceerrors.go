package service

import (
	"errors"
	"fmt"
)

var (
	// ErrNRTM4VersionMismatch nrtm version is not 4
	ErrNRTM4VersionMismatch = errors.New("nrtm version is not 4")
	// ErrNRTM4SourceMismatch session id does not match source
	ErrNRTM4SourceMismatch = errors.New("session id does not match source")
	// ErrNRTM4SourceNameMismatch source name does not match source
	ErrNRTM4SourceNameMismatch = errors.New("source name does not match source")
	// ErrNRTM4NotificationOutOfDate notification file is stale
	ErrNRTM4NotificationOutOfDate = errors.New("notification file is stale")
	// ErrNRTM4FileVersionMismatch file version does not match its reference
	ErrNRTM4FileVersionMismatch = errors.New("file version does not match its reference")
	// ErrNRTM4FileVersionInconsistency version is lower than source
	ErrNRTM4FileVersionInconsistency = errors.New("version is lower than source")
	// ErrNRTM4NoDeltasInNotification the NRTM server published a notification file with no deltas
	ErrNRTM4NoDeltasInNotification = errors.New("no deltas listed in notification file")
	// ErrNRTM4NotificationDeltaSequenceBroken the NRTM server has an incontiguous list of delta version
	ErrNRTM4NotificationDeltaSequenceBroken = errors.New("server has incontiguous list of delta versions")
	// ErrNRTM4NotificationVersionDoesNotMatchDelta the highest delta version is not the notification version
	ErrNRTM4NotificationVersionDoesNotMatchDelta = errors.New("highest delta version is not the notification version")
	// ErrNRTM4DuplicateDeltaVersion the highest delta version is not the notification version
	ErrNRTM4DuplicateDeltaVersion = errors.New("notification file published a duplicate delta file")
)

// ErrNRTMServiceError is when sth is wrong with the NRTM server
type ErrNRTMServiceError struct {
	Message string
}

func (e ErrNRTMServiceError) Error() string {
	return e.Message
}

func newNRTMServiceError(msg string, args ...any) ErrNRTMServiceError {
	return ErrNRTMServiceError{fmt.Sprintf(msg, args...)}
}
