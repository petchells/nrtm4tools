package nrtm4serve

import (
	"log"
	"net/http"
	"time"

	"github.com/petchells/nrtm4tools/internal/nrtm4/pg"
	"github.com/petchells/nrtm4tools/internal/nrtm4/service"
	"github.com/petchells/nrtm4tools/internal/nrtm4serve/rpc"
)

// Launch sets up the rpc handler and starts the server
func Launch(config service.AppConfig, port int, webDir string) {
	repo := pg.PostgresRepository{}
	if err := repo.Initialize(config.PgDatabaseURL); err != nil {
		log.Fatal("Failed to initialize repository")
	}
	defer repo.Close()
	processor := service.NewNRTMProcessor(config, repo, service.HTTPClient{})
	rpcHandler := rpc.Handler{API: WebAPI{Processor: processor}}
	logger.Info("NRTM4serve is starting", "port", port)
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Recovered from Panic in launcher", "recover", r)
			time.Sleep(time.Second * 20)
		}
	}()
	s := rpc.NewServer()
	s.Router().HandleFunc("/rpc", rpcHandler.ProcessRPC).Methods("POST")
	s.Router().HandleFunc("/rpc", rpcHandler.ProcessRPC).Methods("OPTIONS")

	if len(webDir) > 0 {
		s.Router().PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(webDir))))

	}
	s.Serve(port)
}
