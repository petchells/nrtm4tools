package util

import (
	"log/slog"
	"os"
)

// Logger global logger
var Logger = slog.New(
	slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		},
	),
)
