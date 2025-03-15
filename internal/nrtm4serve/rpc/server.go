package rpc

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type route struct {
	pattern *regexp.Regexp
	handler http.Handler
}

// Server needs to know what port to run on
type Server struct {
	r *mux.Router
}

// NewServer creates a new server on `port`
func NewServer() Server {
	r := mux.NewRouter()
	return Server{r}
}

// Serve starts the server
func (s *Server) Serve(port int) error {
	s.r.Walk(
		func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			ht, err := route.GetPathRegexp()
			logger.Debug("Serving route", "path", ht)
			return err
		},
	)
	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return err
	}
	srv := &http.Server{
		Handler:      s.r,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	errc := make(chan error, 1)
	go func() {
		errc <- srv.Serve(l)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errc:
		logger.Error("Failed to serve", "error", err)
	case sig := <-sigs:
		logger.Info("terminating", "sig", sig)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return srv.Shutdown(ctx)
}

// Router returns the router for this server
func (s *Server) Router() *mux.Router {
	return s.r
}
