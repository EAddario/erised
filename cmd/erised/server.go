package main

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

type server struct {
	mux *http.ServeMux
	cfg *http.Server
	ctx context.Context
	stp context.CancelFunc
	pth string
}

func newServer(port, read, write, idle int, path string) *server {
	log.Debug().Msg("entering newServer")
	srv := &server{}
	srv.mux = &http.ServeMux{}

	srv.cfg = &http.Server{
		Addr:         ":" + strconv.Itoa(port),
		Handler:      srv.mux,
		ReadTimeout:  time.Duration(read) * time.Second,
		WriteTimeout: time.Duration(write) * time.Second,
		IdleTimeout:  time.Duration(idle) * time.Second,
	}

	srv.ctx, srv.stp = context.WithCancel(context.Background())
	srv.pth = path
	srv.routes()
	log.Info().
		Str("version", version).
		Int("port", port).
		Str("readTimeout", srv.cfg.ReadTimeout.String()).
		Str("writeTimeout", srv.cfg.WriteTimeout.String()).
		Str("idleTimeout", srv.cfg.IdleTimeout.String()).
		Str("responseFileSearchPath", path).
		Msg("erised server running")
	log.Debug().Msg("leaving newServer")
	return srv
}
