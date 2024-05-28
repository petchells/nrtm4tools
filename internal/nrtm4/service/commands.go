package service

import (
	"regexp"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
)

// AppConfig application configuration object
type AppConfig struct {
	NRTMFilePath     string
	PgDatabaseURL    string
	BoltDatabasePath string
}

// CommandProcessor top-level processing for input commands
type CommandProcessor struct {
	AppConfig
	repo   persist.Repository
	client Client
}

// NewCommandProcessor creates a CommandProcessor
func NewCommandProcessor(config AppConfig, repo persist.Repository, client Client) CommandProcessor {
	return CommandProcessor{
		AppConfig: config,
		repo:      repo,
		client:    client,
	}
}

var labelRegex = regexp.MustCompile("^[A-Za-z][A-Za-z0-9._-]*[A-Za-z]$")

// Connect establishes a new connection to a NRTM source server
func (p CommandProcessor) Connect(notificationURL string, label string) {
	processor := NRTMProcessor{
		config: p.AppConfig,
		repo:   p.repo,
		client: p.client,
	}
	if len(label) > 0 && !labelRegex.MatchString(label) {
		logger.Error("Label must be alphanumeric")
		return
	}
	// TODO
	// Sanitize arguments
	// -- ensure URL looks like a URL, make schema/host lowercase
	err := processor.Connect(notificationURL, label)
	if err != nil {
		logger.Error("Failed to Connect", "url", notificationURL, err)
		return
	}
	logger.Info("Connect successful", "url", notificationURL)
}
