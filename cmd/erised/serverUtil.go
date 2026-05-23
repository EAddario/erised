package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	encodingTEXT = iota
	encodingJSON
	encodingXML
	encodingGZIP
	encodingHTML
)

func elapsedTime(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Debug().Msg(name + " ran for " + elapsed.Round(time.Second).String())
}

func (srv *server) respond(res http.ResponseWriter, encoding int, delay time.Duration, data interface{}) {
	log.Debug().Msg("entering respond")

	if delay > 0 {
		log.Warn().Str("delay", delay.String()).Msg("pausing execution")
		time.Sleep(delay)
	}

	if data == nil {
		data = ""
	}

	switch encoding {
	case encodingTEXT, encodingJSON, encodingXML, encodingHTML:
		if _, err := io.WriteString(res, fmt.Sprintf("%v", data)); err != nil {
			log.Error().Msg(err.Error())
		}
	case encodingGZIP:
		encoder := gzip.NewWriter(res)
		if _, err := encoder.Write([]byte(fmt.Sprintf("%v", data))); err != nil {
			log.Error().Msg(err.Error())
		}
		defer func() {
			if err := encoder.Close(); err != nil {
				log.Error().Msg(err.Error())
			}
		}()
	}

	log.Debug().Msg("leaving respond")
}
