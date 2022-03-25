package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

func (s *server) routes() {
	log.Debug().Msg("entering routes")

	go s.mux.HandleFunc("/", s.handleLanding())
	go s.mux.HandleFunc("/erised/headers", s.handleHeaders())
	go s.mux.HandleFunc("/erised/info", s.handleInfo())
	go s.mux.HandleFunc("/erised/ip", s.handleIP())
	go s.mux.HandleFunc("/erised/shutdown", s.handleShutdown())

	log.Debug().Msg("leaving routes")
}

func (s *server) handleLanding() http.HandlerFunc {
	log.Debug().Msg("entering handleLanding")

	return func(res http.ResponseWriter, req *http.Request) {
		log.Info().
			Str("protocol", req.Proto).
			Str("remoteAddress", req.RemoteAddr).
			Str("method", req.Method).
			Str("host", req.Host).
			Str("uri", req.RequestURI).
			Str("searchPath", *s.pth).
			Msg("handleLanding")

		delay := time.Duration(0)
		enc, ct, ce := encoding(req.Header.Get("X-Erised-Content-Type"))

		res.Header().Set("Content-Type", ct)
		res.Header().Set("Content-Encoding", ce)

		if rd, err := strconv.Atoi(req.Header.Get("X-Erised-Response-Delay")); rd > 0 && err == nil {
			delay = time.Duration(rd) * time.Millisecond
		}

		hd := req.Header.Get("X-Erised-Headers")
		var rs map[string]interface{}

		if err := json.Unmarshal([]byte(hd), &rs); err == nil {
			if len(rs) != 0 {
				for k, v := range rs {
					res.Header().Set(k, fmt.Sprintf("%v", v))
				}
			}
		}

		sc := httpStatusCode(req.Header.Get("X-Erised-Status-Code"))

		if sc >= 300 && sc < 310 {
			res.Header().Set("Location", req.Header.Get("X-Erised-Location"))
		}

		res.WriteHeader(sc)

		data := ""
		if fn := req.Header.Get("X-Erised-Response-File"); fn != "" {

			err := filepath.Walk(*s.pth, func(path string, info os.FileInfo, err error) error {

				if err != nil {
					log.Error().Msg("Invalid path: " + path)
					return nil
				}

				if !info.IsDir() && filepath.Base(path) == fn {
					if ct, err := ioutil.ReadFile(path); err != nil {
						log.Error().Msg("Unable to open the file: " + path)
					} else {
						data = string(ct)
					}
				}

				return nil
			})

			if data == "" || err != nil {
				log.Error().Msg("File not found: " + fn)
			}

		} else {
			data = req.Header.Get("X-Erised-Data")
		}
		s.respond(res, enc, delay, data)

		log.Debug().Msg("leaving handleLanding")
	}
}

func (s *server) handleHeaders() http.HandlerFunc {
	log.Debug().Msg("entering handleHeaders")

	return func(res http.ResponseWriter, req *http.Request) {
		log.Info().
			Str("protocol", req.Proto).
			Str("remoteAddress", req.RemoteAddr).
			Str("method", req.Method).
			Str("host", req.Host).
			Str("uri", req.RequestURI).
			Msg("handleHeaders")

		res.Header().Set("Content-Type", "application/json")
		data := "{"

		for k, v := range req.Header {
			if k == "X-Erised-Data" {
				if json.Valid([]byte(v[0])) {
					data += "\"" + k + "\":" + v[0] + ","
				} else {
					data += "\"" + k + "\":\"" + strings.ReplaceAll(v[0], `"`, `\"`) + "\","
				}
			} else {
				data += "\"" + k + "\":\"" + v[0] + "\","
			}
		}

		data += "\"Host\":\"" + req.Host + "\""
		data += "}"
		s.respond(res, encodingJSON, 0, data)

		log.Debug().Msg("leaving handleHeaders")
	}
}

func (s *server) handleInfo() http.HandlerFunc {
	log.Debug().Msg("entering handleInfo")

	return func(res http.ResponseWriter, req *http.Request) {
		log.Info().
			Str("protocol", req.Proto).
			Str("remoteAddress", req.RemoteAddr).
			Str("method", req.Method).
			Str("host", req.Host).
			Str("uri", req.RequestURI).
			Msg("handleInfo")

		res.Header().Set("Content-Type", "application/json")

		data := "{"
		data += "\"Host\":\"" + req.Host + "\","
		data += "\"Method\":\"" + req.Method + "\","
		data += "\"Protocol\":\"" + req.Proto + "\","
		data += "\"Request URI\":\"" + req.RequestURI + "\""
		data += "}"

		s.respond(res, encodingJSON, 0, data)

		log.Debug().Msg("leaving handleInfo")
	}
}

func (s *server) handleIP() http.HandlerFunc {
	log.Debug().Msg("entering handleIP")

	return func(res http.ResponseWriter, req *http.Request) {
		log.Info().
			Str("protocol", req.Proto).
			Str("remoteAddress", req.RemoteAddr).
			Str("method", req.Method).
			Str("host", req.Host).
			Str("uri", req.RequestURI).
			Msg("handleIP")

		res.Header().Set("Content-Type", "application/json")

		data := "{"
		data += "\"Client IP\":\"" + req.RemoteAddr + "\""
		data += "}"

		s.respond(res, encodingJSON, 0, data)

		log.Debug().Msg("leaving handleIP")
	}
}

func (s *server) handleShutdown() http.HandlerFunc {
	log.Debug().Msg("entering handleShutdown")

	return func(res http.ResponseWriter, req *http.Request) {
		log.Info().
			Str("protocol", req.Proto).
			Str("remoteAddress", req.RemoteAddr).
			Str("method", req.Method).
			Str("host", req.Host).
			Str("uri", req.RequestURI).
			Msg("handleShutdown")

		res.Header().Set("Content-Type", "application/json")

		s.respond(res, encodingJSON, 0, "{\"shutdown\":\"ok\"}")

		if err := s.cfg.Shutdown(context.Background()); err != nil {
			log.Error().Msg(err.Error())
		}

		log.Debug().Msg("leaving handleShutdown")
	}
}
