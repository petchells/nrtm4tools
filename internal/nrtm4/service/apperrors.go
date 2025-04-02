package service

import "errors"

var (
	// File errors

	// ErrHashMismatch when a file downloaded from 'url' does not match its 'hash'
	ErrHashMismatch = errors.New("hash does not match downloaded file")

	// ErrSnapshotInsertFailed snapshot insertion failed
	ErrSnapshotInsertFailed = errors.New("snapshot was not inserted into the repository")

	// Repo errors

	// ErrSessionRestarted server has started a new session
	ErrSessionRestarted = errors.New("server has started a new session")

	// ErrBadNotificationURL notification URL cannot be parsed
	ErrBadNotificationURL = errors.New("notification URL cannot be parsed")

	// ErrSourceAlreadyExists a source with the given label already exists
	ErrSourceAlreadyExists = errors.New("a source with the given label already exists")

	// ErrInvalidLabel label is too lot or contains character which are not allowed
	ErrInvalidLabel = errors.New("label is too long or contains characters which are not allowed")

	// ErrSourceNotFound a source with the given label is not in the repo
	ErrSourceNotFound = errors.New("cannot find source with given name and label")

	// ErrNextConsecutiveDeltaUnavaliable cannot find the next consecutive delta to apply to our repo
	ErrNextConsecutiveDeltaUnavaliable = errors.New("repository is too old to update from the server")
)
