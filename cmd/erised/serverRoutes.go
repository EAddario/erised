package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

func (s *server) routes() {
	log.Debug().Msg("entering routes")

	s.mux.HandleFunc("/", s.handleLanding())
	s.mux.HandleFunc("/erised/headers", s.handleHeaders())
	s.mux.HandleFunc("/erised/ip", s.handleIP())
	s.mux.HandleFunc("/erised/info", s.handleInfo())

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
			Msg("handleLanding")


		delay := time.Duration(0)
		enc, ct, ce := encoding(req.Header.Get("X-Erised-Content-Type"))

		res.Header().Set("Content-Type", ct)
		res.Header().Set("Content-Encoding", ce)

		if rd, err := strconv.Atoi(req.Header.Get("X-Erised-Response-Delay")); err == nil {
			delay = time.Duration(rd) * time.Millisecond
		} else {
			log.Error().Msg(err.Error())
		}

		hd := req.Header.Get("X-Erised-Headers")
		var rs map[string]interface{}

		if err := json.Unmarshal([]byte(hd), &rs); err == nil {
			if len(rs) != 0 {
				for k, v := range rs {
					res.Header().Set(k, fmt.Sprintf("%v", v))
				}
			}
		} else {
			log.Error().Msg(err.Error())
		}

		sc := httpStatusCode(req.Header.Get("X-Erised-Status-Code"))

		if sc >= 300 && sc < 310 {
			res.Header().Set("Location", req.Header.Get("X-Erised-Location"))
		}

		res.WriteHeader(sc)
		data := req.Header.Get("X-Erised-Data")
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
