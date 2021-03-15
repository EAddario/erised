package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const version = "v0.2.4"

type server struct {
	mux *http.ServeMux
	cfg *http.Server
}

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func newServer(port, read, write, idle int) *server {
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
	s.routes()

	log.Log().
		Str("version", version).
		Int("port", port).
		Str("readTimeout", s.cfg.ReadTimeout.String()).
		Str("writeTimeout", s.cfg.WriteTimeout.String()).
		Str("idleTimeout", s.cfg.IdleTimeout.String()).
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
	lv := flag.String("level", "info", "one of debug/warn/error/off. info used otherwise")
	lf := flag.Bool("json", false, "uses JSON log format")

	setupFlags(flag.CommandLine)
	flag.Parse()

	switch strings.ToLower(*lv) {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
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

	srv := newServer(*pt, *rt, *wt, *it)

	if err := srv.cfg.ListenAndServe(); err != nil {
		log.Fatal().Msg(err.Error())
	}

	defer log.Log().Msg("erised server shutting down")
	log.Debug().Msg("leaving main")
}
