package service

import (
	"regexp"
)

// AppConfig application configuration object
type AppConfig struct {
	NRTMFilePath     string
	PgDatabaseURL    string
	BoltDatabasePath string
}

// Processor top-level processing for app functions
type Processor interface {
	Connect(string, string) error
	Update(string, string) error
}

// CommandExecutor top-level processing for input commands
type CommandExecutor struct {
	processor Processor
}

// NewCommandProcessor creates a CommandProcessor
func NewCommandProcessor(processor Processor) CommandExecutor {
	return CommandExecutor{processor}
}

var labelRegex = regexp.MustCompile("^[A-Za-z0-9._-]*[A-Za-z0-9][A-Za-z0-9._-]*$")

// Connect establishes a new connection to a NRTM source server
func (ce CommandExecutor) Connect(notificationURL string, label string) {
	if len(label) > 0 && !labelRegex.MatchString(label) {
		logger.Error("Label must be alphanumeric")
		return
	}
	// TODO
	// Sanitize arguments
	// -- ensure URL looks like a URL, make schema/host lowercase
	err := ce.processor.Connect(notificationURL, label)
	if err != nil {
		logger.Error("Failed to Connect", "url", notificationURL, "error", err)
		return
	}
	logger.Info("Connect successful", "url", notificationURL)
}

// Update brings local mirror up to date
func (ce CommandExecutor) Update(source string, label string) {
	err := ce.processor.Update(source, label)
	if err != nil {
		logger.Warn("Error occurred during update", "error", err)
	} else {
		logger.Info("Update finished successfully")
	}
}
