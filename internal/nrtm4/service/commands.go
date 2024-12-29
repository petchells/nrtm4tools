package service

import (
	"fmt"
	"regexp"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
)

// AppConfig application configuration object
type AppConfig struct {
	NRTMFilePath     string
	PgDatabaseURL    string
	BoltDatabasePath string
}

// ExecutionProcessor top-level processing for app functions
type ExecutionProcessor interface {
	Connect(string, string) error
	Update(string, string) error
	ListSources() ([]persist.NRTMSourceDetails, error)
}

// CommandExecutor top-level processing for input commands
type CommandExecutor struct {
	processor ExecutionProcessor
}

// NewCommandProcessor creates a CommandProcessor
func NewCommandProcessor(processor ExecutionProcessor) CommandExecutor {
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

// ListSources shows all sources in db
func (ce CommandExecutor) ListSources(src, label string) {
	logger.Info("Not doing anything with these args for now", "src", src, "label", label)
	sources, err := ce.processor.ListSources()
	if err != nil {
		logger.Warn("Error occurred when listing sources", "error", err)
		return
	}
	fmt.Printf("   Source                         Label                         \n")
	fmt.Printf("----------------------------------------------------------------\n")
	for i, src := range sources {
		fmt.Printf("%2d %-30v %-30v\n", i+1, src.Source, src.Label)
	}
	logger.Info("List finished successfully")
}

// ReplaceLabel Replaces a label for a source/label
func (ce CommandExecutor) ReplaceLabel(src, fromLabel, toLabel string) {
	logger.Debug("Not doing anything with these args for now", "src", src, "toLabel", toLabel, "fromLabel", fromLabel)
}
