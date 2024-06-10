package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const version = "v0.6.7"

type server struct {
	mux *http.ServeMux
	cfg *http.Server
	pth *string
}

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func newServer(port, read, write, idle int, path string) *server {
	log.Debug().Msg("entering newServer")

	s := &server{}
	s.mux = &http.ServeMux{}
	s.cfg = &http.Server{
		Addr:         ":" + strconv.Itoa(port),
		Handler:      s.mux,
		ReadTimeout:  time.Duration(read) * time.Second,
		WriteTimeout: time.Duration(write) * time.Second,
		IdleTimeout:  time.Duration(idle) * time.Second,
	}
	s.pth = &path
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

func setupFlags(f *flag.FlagSet) {
	log.Debug().Msg("entering setupFlags")

	f.Usage = func() {
		fmt.Println("Simple http server to test arbitrary responses (" + version + ")")
		fmt.Println("Usage examples at https://github.com/EAddario/erised")
		fmt.Println("\nerised [options]")
		fmt.Println("\nParameters:")
		flag.PrintDefaults()
		fmt.Println("\nHTTP Headers:")
		fmt.Println("X-Erised-Content-Type:\t\tSets the response Content-Type")
		fmt.Println("X-Erised-Data:\t\t\tReturns the same value in the response body")
		fmt.Println("X-Erised-Headers:\t\tReturns the value(s) in the response header(s). Values must be in a JSON array")
		fmt.Println("X-Erised-Location:\t\tSets the response Location when 300 ≤ X-Erised-Status-Code < 310")
		fmt.Println("X-Erised-Response-Delay:\tNumber of milliseconds to wait before sending response back to client")
		fmt.Println("X-Erised-Response-File:\t\tReturns the contents of file in the response body. If present, X-Erised-Data is ignored")
		fmt.Println("X-Erised-Status-Code:\t\tSets the HTTP Status Code")
		fmt.Println()
	}

	log.Debug().Msg("leaving setupFlags")
}

func main() {
	log.Debug().Msg("entering main")

	pt := flag.Int("port", 8080, "port to listen")
	rt := flag.Int("read", 5, "maximum duration in seconds for reading the entire request")
	wt := flag.Int("write", 10, "maximum duration in seconds before timing out response writes")
	it := flag.Int("idle", 120, "maximum time in seconds to wait for the next request when keep-alive is enabled")
	lv := flag.String("level", "info", "one of debug/info/warn/error/off")
	lf := flag.Bool("json", false, "use JSON log format")
	ph := flag.String("path", ".", "path to search recursively for X-Erised-Response-File")

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

	srv := newServer(*pt, *rt, *wt, *it, *ph)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		if err := srv.cfg.Close(); err != nil {
			log.Fatal().Msg("HTTP close error: " + err.Error())
		}
	}()

	if err := srv.cfg.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			log.Warn().Msg(err.Error())
		case errors.Is(err, http.ErrBodyReadAfterClose), errors.Is(err, http.ErrHandlerTimeout), errors.Is(err, http.ErrLineTooLong):
			log.Error().Msg(err.Error())
		default:
			log.Fatal().Msg(err.Error())
		}
	}

	defer func() {
		log.Log().Msg("erised server shutting down")
		time.Sleep(100 * time.Millisecond)
	}()

	log.Debug().Msg("leaving main")
}
