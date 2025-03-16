package nrtm4serve

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/petchells/nrtm4tools/internal/nrtm4/pg"
	"github.com/petchells/nrtm4tools/internal/nrtm4/service"
	"github.com/petchells/nrtm4tools/internal/nrtm4/util"
	"github.com/petchells/nrtm4tools/internal/nrtm4serve/rpc"
)

// ClientConfig is read by the web client when it starts
type ClientConfig struct {
	WebSocketURL string
	RPCEndpoint  string
}

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

	s.Router().HandleFunc("/ws", wsHandler)

	serveIndex := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(webDir, "index.html"))
	}

	serveConfig := func(w http.ResponseWriter, r *http.Request) {
		cc := ClientConfig{
			WebSocketURL: config.WebSocketURL,
			RPCEndpoint:  config.RPCEndpoint,
		}
		ccJSON, err := json.Marshal(cc)
		if err != nil {
			log.Fatal("Cannot serialize client config", cc)
		}
		content := bytes.NewReader(ccJSON)
		http.ServeContent(w, r, "webclient.cfg", util.AppClock.Now(), content)
	}
	s.Router().HandleFunc("/s/webclient.cfg", serveConfig).Methods("GET")
	if len(webDir) > 0 {
		s.Router().PathPrefix("/assets/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(webDir))))
		s.Router().HandleFunc("/", serveIndex).Methods("GET")
		s.Router().HandleFunc("/{.*}", serveIndex).Methods("GET")
	}
	log.Fatal(s.Serve(port))
}
