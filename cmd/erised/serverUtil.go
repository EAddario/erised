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

func httpStatusCode(code string) int {
	switch code {
	case "MultipleChoices", "300":
		return 300
	case "MovedPermanently", "301":
		return 301
	case "Found", "302":
		return 302
	case "SeeOther", "303":
		return 303
	case "UseProxy", "305":
		return 305
	case "TemporaryRedirect", "307":
		return 307
	case "PermanentRedirect", "308":
		return 308
	case "BadRequest", "400":
		return 400
	case "Unauthorized", "401":
		return 401
	case "PaymentRequired", "402":
		return 402
	case "Forbidden", "403":
		return 403
	case "NotFound", "404":
		return 404
	case "MethodNotAllowed", "405":
		return 405
	case "RequestTimeout", "408":
		return 408
	case "Conflict", "409":
		return 409
	case "Gone", "410":
		return 410
	case "Teapot", "418":
		return 418
	case "TooManyRequests", "429":
		return 429
	case "InternalServerError", "500":
		return 500
	case "NotImplemented", "501":
		return 501
	case "BadGateway", "502":
		return 502
	case "ServiceUnavailable", "503":
		return 503
	case "GatewayTimeout", "504":
		return 504
	case "HTTPVersionNotSupported", "505":
		return 505
	case "InsufficientStorage", "507":
		return 507
	case "LoopDetected", "508":
		return 508
	case "NotExtended", "510":
		return 510
	case "NetworkAuthenticationRequired", "511":
		return 511
	default:
		return 200
	}
}

func mimeType(code string) (int, string, string) {
	switch code {
	case "json":
		return encodingJSON, "application/json", ""
	case "xml":
		return encodingXML, "application/xml", ""
	case "gzip":
		return encodingGZIP, "application/octet-stream", "gzip"
	case "html":
		return encodingHTML, "text/html", ""
	default:
		return encodingTEXT, "text/plain", ""
	}
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
