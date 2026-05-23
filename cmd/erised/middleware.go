package main

import (
	"compress/gzip"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func WithGzip(next http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Header.Get("X-Erised-Content-Type") == "gzip" {
			res.Header().Set("Content-Encoding", "gzip")
			gz := gzip.NewWriter(res)
			defer gz.Close()
			gzw := &gzipResponseWriter{ResponseWriter: res, Writer: gz}
			next(gzw, req)
		} else {
			next(res, req)
		}
	}
}

func WithDelay(next http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if xrd, err := strconv.Atoi(req.Header.Get("X-Erised-Response-Delay")); xrd > 0 && err == nil {
			delay := time.Duration(xrd) * time.Millisecond
			log.Warn().Str("delay", delay.String()).Msg("pausing execution")
			time.Sleep(delay)
		}
		next(res, req)
	}
}
