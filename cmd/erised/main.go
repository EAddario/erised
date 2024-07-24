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
	"runtime/pprof"
	"strings"
	"syscall"
	"time"
)


func main() {
	defer elapsedTime(time.Now(), "Erised Server")
	log.Debug().Msg("entering main")

	var dir string
	var err error
	idleTimeout := flag.Int("idle", 120, "maximum time in seconds to wait for the next request when keep-alive is enabled")
	jsonLog := flag.Bool("json", false, "use JSON log format")
	logLevel := flag.String("level", "info", "one of debug/info/warn/error/off")
	searchPath := flag.String("path", "", "path to search recursively for X-Erised-Response-File")
	port := flag.Int("port", 8080, "port to listen")
	profile := flag.String("profile", "", "profile this session. A valid file name is required")
	readTimeout := flag.Int("read", 5, "maximum duration in seconds for reading the entire request")
	writeTimeout := flag.Int("write", 10, "maximum duration in seconds before timing out response writes")
	setupFlags(flag.CommandLine)
	flag.Parse()

	if dir, err = os.Getwd(); err != nil {
		panic("Unable to get current directory. Program will terminate.\n\n" + err.Error())
	}

	if *profile != "" {
		if f, err := os.Create(*profile + ".prof"); err == nil {
			if err = pprof.StartCPUProfile(f); err != nil {
				panic("Cannot enable profiling. Program will terminate.\n\n" + err.Error())
			} else {
				defer pprof.StopCPUProfile()
			}
		} else {
			log.Error().Msg("Unable to create profiling file: " + err.Error())
		}
	}

	switch strings.ToLower(*logLevel) {
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

	if *jsonLog {
		log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	}

	if *searchPath != "" {
		*searchPath = filepath.Join(dir, *searchPath)
	}

	srv := newServer(*port, *readTimeout, *writeTimeout, *idleTimeout, *searchPath)

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
