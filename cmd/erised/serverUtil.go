package main

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	encodingTEXT = iota
	encodingJSON
	encodingXML
	encodingBASE64
)

func httpStatusCode(code string) int {
	switch code {
	case "OK", "200" : return 200

	case "MultipleChoices", "300": return 300
	case "MovedPermanently", "301": return 301
	case "Found", "302": return 302
	case "SeeOther", "303": return 303
	case "UseProxy", "305": return 305
	case "TemporaryRedirect", "307": return 307
	case "PermanentRedirect", "308": return 308

	case "BadRequest", "400": return 400
	case "Unauthorized", "401": return 401
	case "PaymentRequired", "402": return 402
	case "Forbidden", "403": return 403
	case "NotFound", "404": return 404
	case "MethodNotAllowed", "405": return 405
	case "RequestTimeout", "408": return 408
	case "Conflict", "409": return 409
	case "Gone", "410": return 410
	case "Teapot", "418": return 418
	case "TooManyRequests", "429": return 429

	case "InternalServerError", "500": return 500
	case "NotImplemented", "501": return 501
	case "BadGateway", "502": return 502
	case "ServiceUnavailable", "503": return 503
	case "GatewayTimeout", "504": return 504
	case "HTTPVersionNotSupported", "505": return 505
	case "InsufficientStorage", "507": return 507
	case "LoopDetected", "508": return 508
	case "NotExtended", "510": return 510
	case "NetworkAuthenticationRequired", "511": return 511
	default: return 200
	}
}

func (s *server) respond(res http.ResponseWriter, encoding int,  data interface{}) {
	if data != nil {

		switch encoding {
		case encodingTEXT:
			if _, err := io.WriteString(res, fmt.Sprintf("%v", data)); err != nil {
				log.Fatal(err)
			}
		case encodingJSON:
			if err := json.NewEncoder(res).Encode(data); err != nil {
				log.Fatal(err)
			}
		case encodingXML:
			if err := xml.NewEncoder(res).Encode(data); err != nil {
				log.Fatal(err)
			}
		case encodingBASE64:
			encoder := base64.NewEncoder(base64.StdEncoding, res)
			if _, err := encoder.Write([]byte(fmt.Sprintf("%v", data))); err != nil {
				log.Fatal(err)
			}
			defer func () {
				if err := encoder.Close(); err != nil {
					log.Fatal(err)
				}
			}()
		}
	}
}
