package cli

import (
	"log"
	"log/slog"
	"os"

	"github.com/petchells/nrtm4tools/internal/nrtm4/pg"
	"github.com/petchells/nrtm4tools/internal/nrtm4/service"
)

// InitializeCommandProcessor starts a db connection pool
func InitializeCommandProcessor(config service.AppConfig) CommandExecutor {
	var httpClient service.HTTPClient
	repo := pg.PostgresRepository{}
	if err := repo.Initialize(config.PgDatabaseURL); err != nil {
		log.Fatal("Failed to initialize repository")
	}
	defer repo.Close()
	service.UserLogger = slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				AddSource: false,
				Level:     slog.LevelInfo,
			},
		),
	)
	processor := service.NewNRTMProcessor(config, repo, httpClient)
	return NewCommandProcessor(processor)
}
