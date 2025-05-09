package service

import (
	"errors"
	"log/slog"
	"net/url"
	"os"
	"strings"
)

// UserLogger is a logger that goes to stdout. Can be overridden.
var UserLogger *slog.Logger

func init() {
	UserLogger = slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				AddSource: false,
				Level:     slog.LevelDebug,
			},
		),
	)
}

func fileNameFromURLString(rawURL string) (string, error) {
	url, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	idx := strings.LastIndex(url.Path, "/")
	if idx > -1 {
		return url.Path[idx+1:], nil
	}
	return "", errors.New("did not find file name in URL")
}
