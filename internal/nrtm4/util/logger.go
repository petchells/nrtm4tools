package util

import (
	"log/slog"
	"os"
)

// Logger global logger
var Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
