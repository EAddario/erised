package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
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
	go s.mux.HandleFunc("/erised/webpage", s.handleWebPage())
	go s.mux.HandleFunc("/erised/webpage/{path...}", s.handleWebPage())
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
			Str("path", req.RequestURI).
			Str("responseFileSearchPath", s.pth).
			Msg("handleLanding")
		delay := time.Duration(0)
		xct := req.Header.Get("X-Erised-Content-Type")
		log.Debug().Msg("X-Erised-Content-Type: " + xct)
		enc, ct, ce := encoding(xct)
		res.Header().Set("Content-Type", ct)

		if xct == "gzip" {
			res.Header().Set("Content-Encoding", ce)
		}

		if rd, err := strconv.Atoi(req.Header.Get("X-Erised-Response-Delay")); rd > 0 && err == nil {
			delay = time.Duration(rd) * time.Millisecond
			log.Debug().Msg("X-Erised-Response-Delay: " + delay.String())
		}

		xhd := req.Header.Get("X-Erised-Headers")
		log.Debug().Msg("X-Erised-Headers: " + xhd)
		var rs map[string]interface{}

		if err := json.Unmarshal([]byte(xhd), &rs); err == nil {
			if len(rs) != 0 {
				for k, v := range rs {
					res.Header().Set(k, fmt.Sprintf("%v", v))
				}
			}
		}

		xsc := httpStatusCode(req.Header.Get("X-Erised-Status-Code"))
		log.Debug().Msg("X-Erised-Status-Code: " + strconv.Itoa(xsc))

		if xsc >= 300 && xsc < 310 {
			xloc := req.Header.Get("X-Erised-Location")
			res.Header().Set("Location", xloc)
			log.Debug().Msg("X-Erised-Location: " + xloc)
		}

		xdt := ""

		if xrf := req.Header.Get("X-Erised-Response-File"); xrf != "" && s.pth != "" {
			log.Debug().Msg("X-Erised-Response-File: " + xrf)
			xsc = http.StatusNotFound

			err := filepath.WalkDir(s.pth, func(path string, entry fs.DirEntry, err error) error {

				if err != nil {
					log.Error().Msg("Invalid path: " + path)
					log.Debug().Msg(fmt.Sprintf("Error: %v", err))

					return errors.New("INVALID_PATH_ERROR")
				}

				if !entry.IsDir() && filepath.Base(path) == xrf {
					if ct, err := os.ReadFile(path); err != nil {
						log.Error().Msg("Unable to open the file: " + path)
						log.Debug().Msg(fmt.Sprintf("Error: %v", err))

						return errors.New("FILE_ACCESS_ERROR")
					} else {
						log.Info().Msg(fmt.Sprintf("Reading file %v", path))
						xdt = string(ct)

						return errors.New("FILE_FOUND")
					}
				}

				log.Debug().Msg("File " + xrf + " not found in " + path)
				return nil
			})

			switch fmt.Sprintf("%v", err) {
			case "INVALID_PATH_ERROR":
				xsc = http.StatusBadRequest
			case "FILE_ACCESS_ERROR":
				xsc = http.StatusInternalServerError
			case "FILE_FOUND":
				xsc = http.StatusOK
			}
		} else {
			xdt = req.Header.Get("X-Erised-Data")
			log.Debug().Msg("X-Erised-Data: " + xdt)
		}

		res.WriteHeader(xsc)
		s.respond(res, enc, delay, xdt)
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
			Str("path", req.RequestURI).
			Msg("handleHeaders")

		if req.Method != http.MethodGet {
			log.Error().Msg("Method " + req.Method + " not allowed for /erised/headers")
			http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

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
			Str("path", req.RequestURI).
			Msg("handleInfo")

		if req.Method != http.MethodGet {
			log.Error().Msg("Method " + req.Method + " not allowed for /erised/info")
			http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

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
			Str("path", req.RequestURI).
			Msg("handleIP")

		if req.Method != http.MethodGet {
			log.Error().Msg("Method " + req.Method + " not allowed for /erised/ip")
			http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

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
			Str("path", req.RequestURI).
			Msg("handleShutdown")

		if req.Method != http.MethodPost {
			log.Error().Msg("Method " + req.Method + " not allowed for /erised/shutdown")
			http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		s.respond(res, encodingJSON, 0, "{\"shutdown\":\"ok\"}")
		log.Info().Msg("Initiating server shutdown")
		s.stp()
		log.Debug().Msg("leaving handleShutdown")
	}
}

func (s *server) handleEchoServer() http.HandlerFunc {
	log.Debug().Msg("entering handleEchoServer")

	return func(res http.ResponseWriter, req *http.Request) {
		log.Info().
			Str("protocol", req.Proto).
			Str("remoteAddress", req.RemoteAddr).
			Str("method", req.Method).
			Str("host", req.Host).
			Str("path", req.RequestURI).
			Msg("handleEchoServer")

		res.Header().Set("Content-Type", "text/html")

		body := ""
		buf := &bytes.Buffer{}
		hn, _ := os.Hostname()

		if _, err := buf.ReadFrom(req.Body); err == nil {
			body = string(buf.Bytes()[:])
		} else {
			log.Error().Msg("Error reading request body")
			log.Debug().Msg(fmt.Sprintf("Error: %v", err))
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		data := "<!DOCTYPE html>"
		data += "<html><head><title>Erised Webpage</title></head>"
		data += "<style>h3 {color: blue; font-family: verdana; margin-bottom: -5px; padding-left: 10px;}"
		data += "p {font-family: courier; margin-bottom: -15px; padding-left: 25px;}</style>"
		data += "<body>"

		se := make([]string, 0, len(os.Environ()))

		for _, env := range os.Environ() {
			se = append(se, env)
		}

		sort.Strings(se)

		data += "<h3><i>Server Environment Variables</i></h3>"
		data += "<p><b>HOSTNAME: </b>" + hn + "</p><br>"

		for _, env := range se {
			ep := strings.SplitN(env, "=", 2)
			k := ep[0]
			v := ep[1]

			data += "<p><b>" + k + ": </b>" + v + "</p>"
		}

		data += "<br><hr><h3><i>Request Info</i></h3>"
		data += "<p><b>Remote Address: </b>" + req.RemoteAddr + "</p>"
		data += "<p><b>Host: </b>" + req.Host + "</p>"
		data += "<p><b>Method: </b>" + req.Method + "</p>"
		data += "<p><b>Protocol: </b>" + req.Proto + "</p>"
		data += "<p><b>Request Path: </b>" + req.RequestURI + "</p>"
		data += "<p><b>Time: </b>" + time.Now().Format(time.RFC850) + "</p>"
		data += "<br><hr><h3><i>Request Headers</i></h3>"

		sh := make([]string, 0, len(req.Header))

		for key := range req.Header {
			sh = append(sh, key)
		}

		sort.Strings(sh)

		for _, k := range sh {
			for _, v := range req.Header[k] {
				data += "<p><b>" + k + ": </b>" + v + "</p>"
			}
		}

		if body != "" {
			data += "<br><hr><h3><i>Request Body</i></h3>"
			data += "<p>" + body + "</p>"
		}

		data += "<br><hr><br><center><a href=\"https://github.com/EAddario/erised\">Erised: A nimble http server to test arbitrary REST API responses.</a></center>"
		data += "</body></html>"
		s.respond(res, encodingHTML, 0, data)
		log.Debug().Msg("leaving handleEchoServer")
	}
}
