package main

import (
	"net/http"
)

func (s *server) routes() {
	s.mux.HandleFunc("/", s.handleLanding())
}

func (s *server) handleLanding() http.HandlerFunc {
	return func (res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/plain; charset=utf-8")

		if ct := req.Header.Get("X-Erised-Content-Type"); ct != "" {
			res.Header().Set("Content-Type", ct)
		}
		if te := req.Header.Get("X-Erised-Transfer-Encoding"); te != "" {
			res.Header().Set("Transfer-Encoding", te)
		}

		res.WriteHeader(httpStatusCode(req.Header.Get("X-Erised-Status-Code")))

		data := req.Header.Get("X-Erised-Data")

		s.respond(res, encodingTEXT, data)
	}
}
