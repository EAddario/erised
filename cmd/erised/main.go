package main

import (
	"context"
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

const version = "v0.11.2"

func main() {
	defer elapsedTime(time.Now(), "Erised Server")
	log.Debug().Msg("entering main")

	var dir string
	var err error
	certFile := flag.String("cert", "", "path to a valid X.509 certificate file")
	idleTimeout := flag.Int("idle", 120, "maximum time in seconds to wait for the next request when keep-alive is enabled")
	jsonLog := flag.Bool("json", false, "use JSON log format")
	keyFile := flag.String("key", "", "path to a valid private key file")
	logLevel := flag.String("level", "info", "one of debug/info/warn/error/off")
	port := flag.Int("port", 0, "port to listen. Default is 8080 for HTTP and 8443 for HTTPS")
	profile := flag.String("profile", "", "profile this session. A valid file name is required")
	readTimeout := flag.Int("read", 5, "maximum duration in seconds for reading the entire request")
	searchPath := flag.String("path", "", "path to search recursively for X-Erised-Response-File")
	useTLS := flag.Bool("https", false, "use HTTPS instead of HTTP. A valid X.509 certificate and private key are required")
	writeTimeout := flag.Int("write", 10, "maximum duration in seconds before timing out response writes")
	setupFlags(flag.CommandLine)
	flag.Parse()

	if *port == 0 && !*useTLS {
		*port = 8080
	} else if *port == 0 && *useTLS {
		*port = 8443
	}

	if dir, err = os.Getwd(); err != nil {
		log.Fatal().Msg("Unable to get current directory. Program will terminate.")
		log.Fatal().Msg(err.Error())
		os.Exit(1)
	}

	if *profile != "" {
		if f, err := os.Create(*profile + ".prof"); err == nil {
			if err = pprof.StartCPUProfile(f); err != nil {
				log.Fatal().Msg("Cannot enable profiling. Program will terminate.")
				log.Fatal().Msg(err.Error())
				os.Exit(1)
			} else {
				defer pprof.StopCPUProfile()
			}
		} else {
			log.Error().Msg("Unable to create profiling file: " + err.Error())
			log.Error().Msg("Profiling will be disabled")
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
	default:
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}

	if *jsonLog {
		log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	}

	if *searchPath != "" {
		*searchPath = filepath.Join(dir, *searchPath)
	}

	if *useTLS && (*certFile == "" || *keyFile == "") {
		log.Fatal().Msg("HTTPS requires a valid certificate and key file")
		os.Exit(1)
	}

	srv := newServer(*port, *readTimeout, *writeTimeout, *idleTimeout, *searchPath)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
		<-sigChan
		srv.stp()
	}()

	go func() {
		if *useTLS {
			if err = srv.cfg.ListenAndServeTLS(*certFile, *keyFile); !errors.Is(err, http.ErrServerClosed) {
				log.Error().Msg("Server shutdown error: " + err.Error())
				if err = syscall.Kill(syscall.Getpid(), syscall.SIGINT); err != nil {
					log.Fatal().Msg(err.Error())
					os.Exit(1)
				}
			}
		} else {
			if err = srv.cfg.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				log.Error().Msg("Server shutdown error: " + err.Error())
				if err = syscall.Kill(syscall.Getpid(), syscall.SIGINT); err != nil {
					log.Fatal().Msg(err.Error())
					os.Exit(1)
				}
			}
		}
	}()

	select {
	case <-srv.ctx.Done():
		if err = srv.cfg.Shutdown(srv.ctx); !errors.Is(err, context.Canceled) {
			if err.Error() != "" {
				log.Error().Msg("Context shutdown error: " + err.Error())
			} else {
				log.Error().Msg("Context shutdown error: unknown reason")
			}
		}
	}

	log.Debug().Msg("leaving main")

	defer func() {
		log.Info().Msg("erised server terminated")
	}()
}
