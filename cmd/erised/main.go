package main

import (
	"errors"
	"flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

const version = "v0.9.7"

func main() {
	cd, err := os.Getwd()

	if err != nil {
		panic("Unable to get current directory. Program will terminate.")
	}

	log.Debug().Msg("entering main")
	pt := flag.Int("port", 8080, "port to listen")
	rt := flag.Int("read", 5, "maximum duration in seconds for reading the entire request")
	wt := flag.Int("write", 10, "maximum duration in seconds before timing out response writes")
	it := flag.Int("idle", 120, "maximum time in seconds to wait for the next request when keep-alive is enabled")
	lv := flag.String("level", "info", "one of debug/info/warn/error/off")
	lf := flag.Bool("json", false, "use JSON log format")
	ph := flag.String("path", "", "path to search recursively for X-Erised-Response-File")
	setupFlags(flag.CommandLine)
	flag.Parse()

	switch strings.ToLower(*lv) {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "off":
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}

	if *lf {
		log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	}

	if *ph != "" {
		*ph = filepath.Join(cd, *ph)
	}

	srv := newServer(*pt, *rt, *wt, *it, *ph)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
		<-sigChan
		srv.stp()
	}()

	go func() {
		if err = srv.cfg.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Error().Msg(err.Error())
			if err = syscall.Kill(syscall.Getpid(), syscall.SIGINT); err != nil {
				panic(err)
			}
		}
	}()

	select {
	case <-srv.ctx.Done():
		if err = srv.cfg.Shutdown(srv.ctx); err != nil {
			log.Error().Msg(err.Error())
		}
	}

	log.Debug().Msg("leaving main")

	defer func() {
		log.Info().Msg("erised server shutting down")
	}()
}
