package cli

import (
	"fmt"

	"github.com/petchells/nrtm4client/internal/nrtm4/persist"
)

// ExecutionProcessor top-level processing for app functions
type ExecutionProcessor interface {
	Connect(string, string) error
	Update(string, string) error
	ListSources() ([]persist.NRTMSourceDetails, error)
	ReplaceLabel(string, string, string) (*persist.NRTMSource, error)
}

// CommandExecutor top-level processing for input commands
type CommandExecutor struct {
	processor ExecutionProcessor
}

// NewCommandProcessor creates a CommandProcessor
func NewCommandProcessor(processor ExecutionProcessor) CommandExecutor {
	return CommandExecutor{processor}
}

// Connect establishes a new connection to a NRTM source server
func (ce CommandExecutor) Connect(notificationURL string, label string) {
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
	logger.Debug("Not doing anything with these args for now", "src", src, "label", label)
	sources, err := ce.processor.ListSources()
	if err != nil {
		logger.Warn("Error occurred when listing sources", "error", err)
		return
	}
	for i, src := range sources {
		fmt.Printf(`		%02d Source    : %v
		Label        : %v
		Version      : %v
		Last updated : %v

`, i+1, src.Source, src.Label, src.Version, src.Notifications[0].Created)
	}
	logger.Info("List finished successfully")
}

// ReplaceLabel Replaces a label for a source/label
func (ce CommandExecutor) ReplaceLabel(src, fromLabel, toLabel string) {
	var updated *persist.NRTMSource
	var err error
	if updated, err = ce.processor.ReplaceLabel(src, fromLabel, toLabel); err != nil {
		logger.Error("ReplaceLabel failed with error", "error", err)
	}
	logger.Info("Replaced label", "updated", updated)
}
