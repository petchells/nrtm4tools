package nrtm4serve

import (
	"bytes"
	"encoding/json"
	"log"
	"log/slog"
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
	Version      string
}

type messagewriter struct {
	hub *Hub
}

func (mw messagewriter) Write(b []byte) (int, error) {
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return 0, err
	}
	msg := message{
		ID:      "logs",
		Content: m,
	}
	mw.hub.send <- msg
	return len(b), nil
}

// Launch sets up the rpc handler and starts the server
func Launch(config service.AppConfig, port int, webDir string) {
	repo := pg.PostgresRepository{}
	if err := repo.Initialize(config.PgDatabaseURL); err != nil {
		log.Fatal("Failed to initialize repository")
	}
	defer repo.Close()
	logger.Info("NRTM4serve is starting", "port", port)
	processor := service.NewNRTMProcessor(config, repo, service.HTTPClient{})
	rpcHandler := rpc.Handler{API: WebAPI{Processor: processor}}
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Recovered from Panic in launcher", "recover", r)
			time.Sleep(time.Second * 20)
		}
	}()
	s := rpc.NewServer()
	s.Router().HandleFunc("/rpc", rpcHandler.ProcessRPC).Methods("POST")
	s.Router().HandleFunc("/rpc", rpcHandler.ProcessRPC).Methods("OPTIONS")

	hub := newHub()
	go hub.run()
	s.Router().HandleFunc("/ws", wsHandler(hub))

	mw := messagewriter{hub}
	service.UserLogger = slog.New(
		slog.NewJSONHandler(
			mw,
			&slog.HandlerOptions{
				AddSource: false,
				Level:     slog.LevelDebug,
			},
		),
	)

	serveIndex := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(webDir, "index.html"))
	}

	cc := ClientConfig{
		WebSocketURL: config.WebSocketURL,
		RPCEndpoint:  config.RPCEndpoint,
		Version:      "alpha",
	}
	ccJSON, err := json.Marshal(cc)
	if err != nil {
		log.Fatal("Cannot serialize client config", cc)
	}
	content := bytes.NewReader(ccJSON)
	serveConfig := func(w http.ResponseWriter, r *http.Request) {
		http.ServeContent(w, r, "clientcfg.json", util.AppClock.Now(), content)
	}
	s.Router().HandleFunc("/s/clientcfg.json", serveConfig).Methods("GET")
	if len(webDir) > 0 {
		s.Router().PathPrefix("/assets/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(webDir))))
		s.Router().HandleFunc("/", serveIndex).Methods("GET")
		s.Router().HandleFunc("/{.*}", serveIndex).Methods("GET")
	}
	log.Fatal(s.Serve(port))
}
