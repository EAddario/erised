package main

import (
	"log"
	"net/http"
	"time"
)

type server struct {
	mux  *http.ServeMux
	cfg  *http.Server
}

func newServer() *server {
	s := &server{}
	s.mux = &http.ServeMux{}
	s.cfg = &http.Server{
		Addr:         ":8080",
		Handler:      s.mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	s.routes()

	return s
}

func main()  {
	srv := newServer()
	log.Fatal(srv.cfg.ListenAndServe())
}
