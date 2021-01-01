package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

const version = "v0.0.3"

type server struct {
	mux *http.ServeMux
	cfg *http.Server
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

	log.Printf("erised %s: Server configured to listen on port %d. Timeouts: %s (read) %s (write) %s (idle)\n",
		version, port, s.cfg.ReadTimeout.String(), s.cfg.WriteTimeout.String(), s.cfg.IdleTimeout.String())

	return s
}

func setupFlags(f *flag.FlagSet) {
	f.Usage = func() {
		fmt.Println("Simple http server to test arbitrary responses (" + version + ")")
		fmt.Println("Usage examples at https://github.com/EAddario/erised")
		fmt.Println("\nerised [options]")
		fmt.Println("\nParameters:")
		flag.PrintDefaults()
		fmt.Println("\nHTTP Headers:")
		fmt.Println("X-Erised-Data:\t\t\tReturns the value in the response body")
		fmt.Println("X-Erised-Content-Type:\t\tReturns the value in the Content-Type response header")
		fmt.Println("X-Erised-Status-Code:\t\tUsed to set the http status code value")
		fmt.Println("X-Erised-Location:\t\tReturns the value of a new URL or path when 300 ≤ X-Erised-Status-Code < 310")
		fmt.Println("X-Erised-Response-Delay:\tNumber of milliseconds to wait before sending response back to client")
		fmt.Println()
	}
}

func main() {
	pt := flag.Int("port", 8080, "port to listen")
	rt := flag.Int("read", 5, "maximum duration in seconds for reading the entire request")
	wt := flag.Int("write", 10, "maximum duration in seconds before timing out response writes")
	it := flag.Int("idle", 120, "maximum time in seconds to wait for the next request when keep-alive is enabled")

	setupFlags(flag.CommandLine)
	flag.Parse()

	srv := newServer(*pt, *rt, *wt, *it)
	log.Fatal(srv.cfg.ListenAndServe())
}
