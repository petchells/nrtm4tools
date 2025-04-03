package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/petchells/nrtm4tools/internal/nrtm4/service"
	"github.com/petchells/nrtm4tools/internal/nrtm4serve"
)

var port = flag.Int("port", 8080, "server port number")
var webdir = flag.String("webdir", "", "(optional) directory containing static web files")
var wsURL = flag.String("wsurl", "", "web socket URL, defaults to http://localhost:<port>/ws")
var rpcURL = flag.String("rpcurl", "", "JSON RPC endpoint URL, defaults to http://localhost:<port>/rpc")

func main() {
	flag.Parse()
	envVars := []string{"PG_DATABASE_URL", "NRTM4_FILE_PATH"}
	for _, ev := range envVars {
		if len(os.Getenv(ev)) <= 0 {
			log.Fatalln("Environment variable not set: ", ev)
		}
	}
	if wsURL == nil || len(*wsURL) == 0 {
		wsu := os.Getenv("WEB_SOCKET_URL")
		if len(wsu) > 0 {
			wsURL = &wsu
		} else {
			wsu = fmt.Sprintf("http://localhost:%d/ws", *port)
			wsURL = &wsu
		}
	}
	if rpcURL == nil || len(*rpcURL) == 0 {
		ru := os.Getenv("JSON_RPC_ENDPOINT")
		if len(ru) > 0 {
			rpcURL = &ru
		} else {
			ru := fmt.Sprintf("http://localhost:%d/rpc", *port)
			rpcURL = &ru
		}
	}
	dbURL := os.Getenv("PG_DATABASE_URL")
	boltDBPath := os.Getenv("BOLT_DATABASE_PATH")
	nrtmFilePath := os.Getenv("NRTM4_FILE_PATH")
	config := service.AppConfig{
		NRTMFilePath:     nrtmFilePath,
		PgDatabaseURL:    dbURL,
		BoltDatabasePath: boltDBPath,
		WebSocketURL:     *wsURL,
		RPCEndpoint:      *rpcURL,
	}
	nrtm4serve.Launch(config, *port, *webdir)
}
