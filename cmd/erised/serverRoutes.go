package main

import (
	"log"
	"net/http"
)

func (s *server) routes() {
	s.mux.HandleFunc("/", s.handleLanding())
}

func (s *server) handleLanding() http.HandlerFunc {
	return func (res http.ResponseWriter, req *http.Request) {
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
