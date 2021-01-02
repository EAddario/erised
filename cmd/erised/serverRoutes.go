package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (s *server) routes() {
	s.mux.HandleFunc("/", s.handleLanding())
	s.mux.HandleFunc("/erised/headers", s.handleHeaders())
	s.mux.HandleFunc("/erised/ip", s.handleIP())
	s.mux.HandleFunc("/erised/info", s.handleInfo())
}

func (s *server) handleLanding() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		log.Printf("%s from %s - %s %s%s",
			req.Proto, req.RemoteAddr, req.Method, req.Host, req.RequestURI)

		res.Header().Set("Content-Type", "text/plain; charset=utf-8")
		delay := time.Duration(0)

		if rd, err := strconv.Atoi(req.Header.Get("X-Erised-Response-Delay")); err == nil {
			delay = time.Duration(rd) * time.Millisecond
		}
		if ct := req.Header.Get("X-Erised-Content-Type"); ct != "" {
			res.Header().Set("Content-Type", ct)
		}
		if te := req.Header.Get("X-Erised-Transfer-Encoding"); te != "" {
			res.Header().Set("Transfer-Encoding", te)
		}
		sc := httpStatusCode(req.Header.Get("X-Erised-Status-Code"))
		if sc >= 300 && sc < 310 {
			res.Header().Set("Location", req.Header.Get("X-Erised-Location"))
		}

		res.WriteHeader(sc)

		data := req.Header.Get("X-Erised-Data")

		s.respond(res, encodingTEXT, delay, data)
	}
}

func (s *server) handleHeaders() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		log.Printf("%s from %s - %s %s%s",
			req.Proto, req.RemoteAddr, req.Method, req.Host, req.RequestURI)

		res.Header().Set("Content-Type", "application/json; charset=utf-8")

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

		s.respond(res, encodingTEXT, 0, data)
	}
}

func (s *server) handleIP() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		log.Printf("%s from %s - %s %s%s",
			req.Proto, req.RemoteAddr, req.Method, req.Host, req.RequestURI)

		res.Header().Set("Content-Type", "application/json; charset=utf-8")

		data := "{"
		data += "\"Client IP\":\"" + req.RemoteAddr + "\""
		data += "}"

		s.respond(res, encodingTEXT, 0, data)
	}
}

func (s *server) handleInfo() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		log.Printf("%s from %s - %s %s%s",
			req.Proto, req.RemoteAddr, req.Method, req.Host, req.RequestURI)

		res.Header().Set("Content-Type", "application/json; charset=utf-8")

		data := "{"
		data += "\"Host\":\"" + req.Host + "\","
		data += "\"Method\":\"" + req.Method + "\","
		data += "\"Protocol\":\"" + req.Proto + "\","
		data += "\"Request URI\":\"" + req.RequestURI + "\""
		data += "}"

		s.respond(res, encodingTEXT, 0, data)
	}
}
