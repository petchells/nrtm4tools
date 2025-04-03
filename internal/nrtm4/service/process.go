package service

import (
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/petchells/nrtm4tools/internal/nrtm4/persist"
)

var (
	fileWriteBufferLength = 1024 * 8
	rpslInsertBatchSize   = 500
)

// AppConfig application configuration object
type AppConfig struct {
	NRTMFilePath     string
	PgDatabaseURL    string
	BoltDatabasePath string
	WebSocketURL     string
	RPCEndpoint      string
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
	UserLogger.Info("Connect to source", "url", notificationURL, "label", label)
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
	fm := fileManager{p.client}
	notification, err := fm.downloadNotificationFile(unfURL)
	if err != nil {
		return err
	}
	err = fm.ensureDirectoryExists(p.config.NRTMFilePath)
	if err != nil {
		return err
	}
	source := persist.NewNRTMSource(notification, label, unfURL)
	source.Status = "new"
	if source, err = ds.saveNewSource(source, notification); err != nil {
		UserLogger.Error("There was a problem saving the source. Remove it and restart sync", "error", err)
		return err
	}
	UserLogger.Info("Saved source", "source", notification.Source, "version", notification.Version, "label", label)
	dirname := filepath.Join(p.config.NRTMFilePath, source.Source, source.SessionID)
	_, err = os.Stat(dirname)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(dirname, 0755); err != nil {
			return err
		}
	}
	// Download snapshot
	UserLogger.Info("Fetching snapshot file", "url", notification.SnapshotRef.URL)
	snapshotFile, err := fm.fetchFileAndCheckHash(unfURL, notification.SnapshotRef, dirname)
	if err != nil {
		source.Status = "snapshot.file.failed: " + err.Error()
		ds.updateSource(source)
		return err
	}
	defer snapshotFile.Close()

	UserLogger.Info("Inserting snapshot objects", "source", notification.Source)
	if err := fm.readJSONSeqRecords(snapshotFile, snapshotObjectInsertFunc(p.repo, source, notification)); err != io.EOF {
		UserLogger.Error("Invalid snapshot. Remove Source and restart sync", "error", err)
		source.Status = "snapshot.insert.failed: " + err.Error()
		ds.updateSource(source)
		return err
	}
	source.Status = "updating"
	ds.updateSource(source)

	UserLogger.Info("Synchronizing deltas", "total refs", len(notification.DeltaRefs))
	source, err = syncDeltas(p, notification, source)
	if err != nil {
		UserLogger.Error("Failed to sync deltas", "source", source.Source, "version", source.Version, "error", err)
		source.Status = "delta.failed: " + err.Error()
		ds.updateSource(source)
		return err
	}
	source.Status = "ok"
	_, err = ds.updateSource(source)
	return err
}

// Update brings the local mirror up to date
func (p NRTMProcessor) Update(sourceName string, label string) (*persist.NRTMSource, error) {
	ds := NrtmDataService{Repository: p.repo}
	source := ds.getSourceByNameAndLabel(sourceName, label)
	if source == nil {
		logger.Warn("No source with given name and label", "name", sourceName, "label", label)
		return nil, ErrSourceNotFound
	}
	fm := fileManager{p.client}
	notification, err := fm.downloadNotificationFile(source.NotificationURL)
	if err != nil {
		return nil, err
	}
	if notification.SessionID != source.SessionID {
		source.Status = "session.restarted"
		ds.updateSource(*source)
		return nil, ErrSessionRestarted
	}
	if notification.Version < int64(source.Version) {
		return nil, ErrNRTM4NotificationOutOfDate
	}
	if notification.Version == int64(source.Version) {
		UserLogger.Info("Already at latest version")
		return nil, nil
	}
	source.Status = "updating"
	source, err = ds.updateSource(*source)
	if err != nil {
		logger.Error("Failed to save source")
	}
	var updated persist.NRTMSource
	if updated, err = syncDeltas(p, notification, *source); err != nil {
		updated.Status = "delta.failed: " + err.Error()
		_, err = ds.updateSource(updated)
		return nil, err
	}
	updated.Status = "ok"
	return ds.updateSource(updated)
}

// ListSources gets details, including notifications, of all sources
func (p NRTMProcessor) ListSources() ([]persist.NRTMSourceDetails, error) {
	ds := NrtmDataService{Repository: p.repo}
	UserLogger.Info("List sources")
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
	UserLogger.Info("Replace label", "src", src, "label", fromLabel, "replace with", toLabel)
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
	return ds.updateSource(*target)
}

// RemoveSource removes a source from the repo
func (p NRTMProcessor) RemoveSource(src, label string) error {
	UserLogger.Info("Remove source", "src", src, "label", label)
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
	sepIdx := 0
	if strings.HasPrefix(relpath, "/") {
		sepIdx = 1
	}
	return base[:idx+1] + relpath[sepIdx:]
}

func validateURLString(str string) bool {
	url, err := url.Parse(str)
	return err == nil && (url.Scheme == "http" || url.Scheme == "https")
}
