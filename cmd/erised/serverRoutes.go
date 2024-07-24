package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

func (srv *server) routes() {
	log.Debug().Msg("entering routes")
	go srv.mux.HandleFunc("/", srv.handleLanding())
	go srv.mux.HandleFunc("/erised/headers", srv.handleHeaders())
	go srv.mux.HandleFunc("/erised/info", srv.handleInfo())
	go srv.mux.HandleFunc("/erised/ip", srv.handleIP())
	go srv.mux.HandleFunc("/erised/shutdown", srv.handleShutdown())
	go srv.mux.HandleFunc("/erised/echoserver", srv.handleEchoServer())
	go srv.mux.HandleFunc("/erised/echoserver/{path...}", srv.handleEchoServer())
	log.Debug().Msg("leaving routes")
}

func (srv *server) handleLanding() http.HandlerFunc {
	log.Debug().Msg("entering handleLanding")

	return func(res http.ResponseWriter, req *http.Request) {
		log.Info().
			Str("protocol", req.Proto).
			Str("remoteAddress", req.RemoteAddr).
			Str("method", req.Method).
			Str("host", req.Host).
			Str("path", req.RequestURI).
			Str("responseFileSearchPath", srv.pth).
			Msg("handleLanding")
		delay := time.Duration(0)
		xContentType := req.Header.Get("X-Erised-Content-Type")
		log.Debug().Msg("X-Erised-Content-Type: " + xContentType)
		encoding, mime, contentEncoding := mimeType(xContentType)
		res.Header().Set("Content-Type", mime)

		if xContentType == "gzip" {
			res.Header().Set("Content-Encoding", contentEncoding)
		}

		if xrd, err := strconv.Atoi(req.Header.Get("X-Erised-Response-Delay")); xrd > 0 && err == nil {
			delay = time.Duration(xrd) * time.Millisecond
			log.Debug().Msg("X-Erised-Response-Delay: " + delay.String())
		}

		xHeaders := req.Header.Get("X-Erised-Headers")
		log.Debug().Msg("X-Erised-Headers: " + xHeaders)
		var hdrs map[string]interface{}

		if err := json.Unmarshal([]byte(xHeaders), &hdrs); err == nil {
			if len(hdrs) != 0 {
				for k, v := range hdrs {
					res.Header().Set(k, fmt.Sprintf("%v", v))
				}
			}
		}

		xStatusCode := httpStatusCode(req.Header.Get("X-Erised-Status-Code"))
		log.Debug().Msg("X-Erised-Status-Code: " + strconv.Itoa(xStatusCode))

		if xStatusCode >= 300 && xStatusCode < 310 {
			xloc := req.Header.Get("X-Erised-Location")
			res.Header().Set("Location", xloc)
			log.Debug().Msg("X-Erised-Location: " + xloc)
		}

		xData := ""

		if xResponseFile := req.Header.Get("X-Erised-Response-File"); xResponseFile != "" && srv.pth != "" {
			log.Debug().Msg("X-Erised-Response-File: " + xResponseFile)
			xStatusCode = http.StatusNotFound

			err := filepath.WalkDir(srv.pth, func(path string, entry fs.DirEntry, err error) error {

				if err != nil {
					log.Error().Msg("Invalid path: " + path)
					log.Debug().Msg(fmt.Sprintf("Error: %v", err))

					return errors.New("INVALID_PATH_ERROR")
				}

				if !entry.IsDir() && filepath.Base(path) == xResponseFile {
					if ct, err := os.ReadFile(path); err != nil {
						log.Error().Msg("Unable to open the file: " + path)
						log.Debug().Msg(fmt.Sprintf("Error: %v", err))

						return errors.New("FILE_ACCESS_ERROR")
					} else {
						log.Info().Msg(fmt.Sprintf("Reading file %v", path))
						xData = string(ct)

						return errors.New("FILE_FOUND")
					}
				}

				log.Debug().Msg("File " + xResponseFile + " not found in " + path)
				return nil
			})

			switch fmt.Sprintf("%v", err) {
			case "INVALID_PATH_ERROR":
				xStatusCode = http.StatusBadRequest
			case "FILE_ACCESS_ERROR":
				xStatusCode = http.StatusInternalServerError
			case "FILE_FOUND":
				xStatusCode = http.StatusOK
			}
		} else {
			xData = req.Header.Get("X-Erised-Data")
			log.Debug().Msg("X-Erised-Data: " + xData)
		}

		res.WriteHeader(xStatusCode)
		srv.respond(res, encoding, delay, xData)
		log.Debug().Msg("leaving handleLanding")
	}
}

func (srv *server) handleHeaders() http.HandlerFunc {
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
		srv.respond(res, encodingJSON, 0, data)
		log.Debug().Msg("leaving handleHeaders")
	}
}

func (srv *server) handleInfo() http.HandlerFunc {
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
		srv.respond(res, encodingJSON, 0, data)
		log.Debug().Msg("leaving handleInfo")
	}
}

func (srv *server) handleIP() http.HandlerFunc {
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
		srv.respond(res, encodingJSON, 0, data)
		log.Debug().Msg("leaving handleIP")
	}
}

func (srv *server) handleShutdown() http.HandlerFunc {
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
		srv.respond(res, encodingJSON, 0, "{\"shutdown\":\"ok\"}")
		log.Info().Msg("Initiating server shutdown")
		srv.stp()
		log.Debug().Msg("leaving handleShutdown")
	}
}

func (srv *server) handleEchoServer() http.HandlerFunc {
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
		hostName, _ := os.Hostname()

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

		env := make([]string, 0, len(os.Environ()))

		for _, v := range os.Environ() {
			env = append(env, v)
		}

		sort.Strings(env)

		data += "<h3><i>Server Environment Variables</i></h3>"
		data += "<p><b>HOSTNAME: </b>" + hostName + "</p><br>"

		for _, v := range env {
			pq := strings.SplitN(v, "=", 2)
			p := pq[0]
			q := pq[1]

			data += "<p><b>" + p + ": </b>" + q + "</p>"
		}

		data += "<br><hr><h3><i>Request Info</i></h3>"
		data += "<p><b>Remote Address: </b>" + req.RemoteAddr + "</p>"
		data += "<p><b>Host: </b>" + req.Host + "</p>"
		data += "<p><b>Method: </b>" + req.Method + "</p>"
		data += "<p><b>Protocol: </b>" + req.Proto + "</p>"
		data += "<p><b>Request Path: </b>" + req.RequestURI + "</p>"
		data += "<p><b>Time: </b>" + time.Now().Format(time.RFC850) + "</p>"
		data += "<br><hr><h3><i>Request Headers</i></h3>"

		hdrs := make([]string, 0, len(req.Header))

		for key := range req.Header {
			hdrs = append(hdrs, key)
		}

		sort.Strings(hdrs)

		for _, k := range hdrs {
			for _, v := range req.Header[k] {
				data += "<p><b>" + k + ": </b>" + v + "</p>"
			}
		}

		if body != "" {
			data += "<br><hr><h3><i>Request Body</i></h3>"
			data += "<p>" + body + "</p>"
		}

		data += "<br><hr><br><center><a href=\"https://github.com/EAddario/erised\">Erised (" + version + "): A nimble http server to test arbitrary REST API responses.</a></center>"
		data += "</body></html>"
		srv.respond(res, encodingHTML, 0, data)
		log.Debug().Msg("leaving handleEchoServer")
	}
}
