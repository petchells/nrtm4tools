package service

import (
	"errors"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/petchells/nrtm4tools/internal/nrtm4/persist"
	"github.com/petchells/nrtm4tools/internal/nrtm4/pg/db"
)

var (
	// Repo errors

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

	fileWriteBufferLength = 1024 * 8
	rpslInsertBatchSize   = 1000
)

// AppConfig application configuration object
type AppConfig struct {
	NRTMFilePath     string
	PgDatabaseURL    string
	BoltDatabasePath string
}

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

const charsAllowedInLabel = "A-Za-z0-9 :._-"

// Must have a letter or digit in there somewhere
var labelRe = regexp.MustCompile("^[" + charsAllowedInLabel + "]*[A-Za-z0-9][" + charsAllowedInLabel + "]*$")

// Connect stores details about a connection
func (p NRTMProcessor) Connect(notificationURL string, label string) error {
	unfURL := strings.TrimSpace(notificationURL)
	if !validateURLString(unfURL) {
		return ErrBadNotificationURL
	}
	label = strings.TrimSpace(label)
	if len(label) > 0 && !labelRe.MatchString(label) {
		return ErrInvalidLabel
	}
	if len(label) > 255 {
		return ErrInvalidLabel
	}
	ds := NrtmDataService{Repository: p.repo}
	if ds.getSourceByURLAndLabel(unfURL, label) != nil {
		return ErrSourceAlreadyExists
	}
	logger.Info("Fetching notification")
	fm := fileManager{p.client}
	notification, err := fm.downloadNotificationFile(unfURL)
	if err != nil {
		return err
	}
	err = fm.ensureDirectoryExists(p.config.NRTMFilePath)
	if err != nil {
		return err
	}
	logger.Info("Saving new source", "source", notification.Source)
	source := persist.NewNRTMSource(notification, label, unfURL)
	if source, err = ds.saveNewSource(source, notification); err != nil {
		logger.Error("There was a problem saving the source. Remove it and restart sync", "error", err)
		return err
	}
	dirname := filepath.Join(p.config.NRTMFilePath, source.Source, source.SessionID)
	_, err = os.Stat(dirname)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(dirname, 0755); err != nil {
			return err
		}
	}
	// Download snapshot
	logger.Info("Fetching snapshot file...")
	snapshotFile, err := fm.fetchFileAndCheckHash(unfURL, notification.SnapshotRef, dirname)
	if err != nil {
		return err
	}
	logger.Info("Snapshot file downloaded")
	defer snapshotFile.Close()

	logger.Info("Inserting snapshot objects", "source", notification.Source)
	if err := fm.readJSONSeqRecords(snapshotFile, snapshotObjectInsertFunc(p.repo, source, notification)); err != io.EOF {
		logger.Error("Invalid snapshot. Remove Source and restart sync", "error", err)
		return err
	}
	return syncDeltas(p, notification, source)
}

// Update brings the local mirror up to date
func (p NRTMProcessor) Update(sourceName string, label string) error {
	ds := NrtmDataService{Repository: p.repo}
	source := ds.getSourceByNameAndLabel(sourceName, label)
	if source == nil {
		logger.Warn("No source with given name and label", "name", sourceName, "label", label)
		return ErrSourceNotFound
	}
	fm := fileManager{p.client}
	notification, err := fm.downloadNotificationFile(source.NotificationURL)
	if err != nil {
		return err
	}
	if notification.SessionID != source.SessionID {
		return errors.New("server has a new mirror session")
	}
	if notification.Version < int64(source.Version) {
		return errors.New("server has old version")
	}
	if notification.Version == int64(source.Version) {
		logger.Info("Already at latest version")
		return nil
	}
	return syncDeltas(p, notification, *source)
}

// ListSources gets details, including notifications, of all sources
func (p NRTMProcessor) ListSources() ([]persist.NRTMSourceDetails, error) {
	ds := NrtmDataService{Repository: p.repo}
	sources, err := ds.listSources()
	deets := []persist.NRTMSourceDetails{}
	if err != nil {
		return deets, err
	}
	for _, src := range sources {
		to := src.Version
		from := src.Version - 99
		if src.Version <= 99 {
			from = 1
		}
		notifs, err := p.repo.GetNotificationHistory(src, from, to)
		if err != nil {
			return deets, err
		}
		deets = append(deets, persist.NRTMSourceDetails{NRTMSource: src, Notifications: notifs})
	}
	return deets, nil
}

// ReplaceLabel replaces a label name
func (p NRTMProcessor) ReplaceLabel(src, fromLabel, toLabel string) (*persist.NRTMSource, error) {
	ds := NrtmDataService{Repository: p.repo}
	possDupe := ds.getSourceByNameAndLabel(src, toLabel)
	if possDupe != nil {
		return nil, ErrSourceAlreadyExists
	}
	target := ds.getSourceByNameAndLabel(src, fromLabel)
	if target == nil {
		return nil, ErrSourceNotFound
	}
	target.Label = toLabel
	return target, db.WithTransaction(func(tx pgx.Tx) error {
		return db.Update(tx, target)
	})
}

// RemoveSource removes a source from the repo
func (p NRTMProcessor) RemoveSource(src, label string) error {
	ds := NrtmDataService{Repository: p.repo}
	target := ds.getSourceByNameAndLabel(src, label)
	if target == nil {
		return ErrSourceNotFound
	}
	return ds.deleteSource(*target)
}

func fullURL(base, relpath string) string {
	idx := strings.LastIndex(base, "/")
	if idx < 0 {
		logger.Error("fullURL called with a path that does not contain '/'", "base", base)
		return ""
	}
	sepIdx := 1
	if strings.HasPrefix(relpath, "/") {
		sepIdx = 0
	}
	return base[:idx+sepIdx] + relpath
}

func validateURLString(str string) bool {
	url, err := url.Parse(str)
	return err == nil && (url.Scheme == "http" || url.Scheme == "https")
}
