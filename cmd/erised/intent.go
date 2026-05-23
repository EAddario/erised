package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrInvalidPath  = errors.New("invalid path error")
	ErrFileAccess   = errors.New("file access error")
	ErrFileNotFound = errors.New("file not found")
)

type FileResolver func(filename string) ([]byte, error)

type ErisedIntent struct {
	StatusCode int
	Headers    map[string]string
	Data       string
}

func BuildErisedIntent(reqHeaders http.Header, resolveFile FileResolver) ErisedIntent {
	intent := ErisedIntent{
		Headers: make(map[string]string),
	}

	// Content Type
	xContentType := reqHeaders.Get("X-Erised-Content-Type")
	mime, contentEncoding := getMimeType(xContentType)
	intent.Headers["Content-Type"] = mime
	if contentEncoding != "" {
		intent.Headers["Content-Encoding"] = contentEncoding
	}

	// Custom Headers
	xHeaders := reqHeaders.Get("X-Erised-Headers")
	var hdrs map[string]interface{}
	if err := json.Unmarshal([]byte(xHeaders), &hdrs); err == nil {
		if len(hdrs) != 0 {
			for k, v := range hdrs {
				intent.Headers[k] = fmt.Sprintf("%v", v)
			}
		}
	}

	// Status Code
	xStatusCode := getHttpStatusCode(reqHeaders.Get("X-Erised-Status-Code"))
	intent.StatusCode = xStatusCode

	// Location
	if xStatusCode >= 300 && xStatusCode < 310 {
		xloc := reqHeaders.Get("X-Erised-Location")
		intent.Headers["Location"] = xloc
	}

	// Response File
	xResponseFile := reqHeaders.Get("X-Erised-Response-File")
	if xResponseFile != "" && resolveFile != nil {
		fileBytes, err := resolveFile(xResponseFile)
		if err != nil {
			switch {
			case errors.Is(err, ErrInvalidPath):
				intent.StatusCode = http.StatusBadRequest
			case errors.Is(err, ErrFileAccess):
				intent.StatusCode = http.StatusInternalServerError
			case errors.Is(err, ErrFileNotFound):
				intent.StatusCode = http.StatusNotFound
			default:
				intent.StatusCode = http.StatusInternalServerError
			}
		} else {
			intent.Data = string(fileBytes)
			intent.StatusCode = http.StatusOK
		}
	} else {
		intent.Data = reqHeaders.Get("X-Erised-Data")
	}

	return intent
}

func getHttpStatusCode(code string) int {
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

func getMimeType(code string) (string, string) {
	switch code {
	case "json":
		return "application/json", ""
	case "xml":
		return "application/xml", ""
	case "gzip":
		return "application/octet-stream", "gzip"
	case "html":
		return "text/html", ""
	default:
		return "text/plain", ""
	}
}
