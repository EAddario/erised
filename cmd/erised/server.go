package main

import (
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
	"time"
)

type server struct {
	mux *mux.Router
	cfg *http.Server
	pth string
}

func newServer(port, read, write, idle int, path string) *server {
	log.Debug().Msg("entering newServer")

	s := &server{}
	s.mux = mux.NewRouter()
	s.cfg = &http.Server{
		Addr:         ":" + strconv.Itoa(port),
		Handler:      s.mux,
		ReadTimeout:  time.Duration(read) * time.Second,
		WriteTimeout: time.Duration(write) * time.Second,
		IdleTimeout:  time.Duration(idle) * time.Second,
	}
	s.pth = path
	s.routes()

	log.Log().
		Str("version", version).
		Int("port", port).
		Str("readTimeout", s.cfg.ReadTimeout.String()).
		Str("writeTimeout", s.cfg.WriteTimeout.String()).
		Str("idleTimeout", s.cfg.IdleTimeout.String()).
		Str("path", path).
		Msg("erised server running")

	log.Debug().Msg("leaving newServer")
	return s
}
