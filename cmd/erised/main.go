package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type server struct {
	mux  *http.ServeMux
	cfg  *http.Server
}

func newServer(port, read, write, idle int) *server {
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

	log.Printf("Server configured to listen on port %s. Timeouts: %s (read) %s (write) %s (idle)\n",
		s.cfg.Addr, s.cfg.ReadTimeout.String(), s.cfg.WriteTimeout.String(), s.cfg.IdleTimeout.String())

	return s
}

func setupFlags(f *flag.FlagSet) {
	f.Usage = func() {
		_, _ = fmt.Printf("%s: a simple http server to test arbitrary responses. Usage example at https://github.com/EAddario/erised\n", filepath.Base(os.Args[0]))
		fmt.Println("\nParameters:")
		flag.PrintDefaults()
	}
}

func main()  {
	pt := flag.Int("port", 8080, "port to listen")
	rt := flag.Int("read", 5, "read timeout in seconds")
	wt := flag.Int("write", 10, "write timeout in seconds")
	it := flag.Int("idle", 120, "idle timeout in seconds")

	setupFlags(flag.CommandLine)
	flag.Parse()

	srv := newServer(*pt, *rt, *wt, *it)
	log.Fatal(srv.cfg.ListenAndServe())
}
