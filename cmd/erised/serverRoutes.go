package main

import (
	"encoding/json"
	"log"
	"net/http"
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

		s.respond(res, encodingTEXT, data)
	}
}

func (s *server) handleHeaders() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		log.Printf("%s from %s - %s %s%s",
			req.Proto, req.RemoteAddr, req.Method, req.Host, req.RequestURI)

		res.Header().Set("Content-Type", "application/json; charset=utf-8")

		if rh, err := json.Marshal(req.Header); err == nil {
			s.respond(res, encodingTEXT, string(rh))
		} else {
			log.Fatal(err)
		}
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

		s.respond(res, encodingTEXT, data)
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

		s.respond(res, encodingTEXT, data)
	}
}
