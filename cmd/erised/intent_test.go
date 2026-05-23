package main

import (
	"net/http"
	"testing"
)

func TestBuildErisedIntent(t *testing.T) {
	t.Run("basic properties (content type, data, status code)", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("X-Erised-Content-Type", "json")
		headers.Set("X-Erised-Data", `{"foo":"bar"}`)
		headers.Set("X-Erised-Status-Code", "Created") // Invalid text, defaults to 200

		intent := BuildErisedIntent(headers, nil)

		if intent.Headers["Content-Type"] != "application/json" {
			t.Errorf("expected json content type, got %v", intent.Headers["Content-Type"])
		}
		if intent.Data != `{"foo":"bar"}` {
			t.Errorf("expected data, got %v", intent.Data)
		}
		if intent.StatusCode != 200 {
			t.Errorf("expected status code 200, got %v", intent.StatusCode)
		}
	})

	t.Run("status code mappings", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("X-Erised-Status-Code", "NotFound")

		intent := BuildErisedIntent(headers, nil)

		if intent.StatusCode != 404 {
			t.Errorf("expected status code 404, got %v", intent.StatusCode)
		}
	})

	t.Run("custom headers", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("X-Erised-Headers", `{"X-Custom":"value"}`)

		intent := BuildErisedIntent(headers, nil)

		if intent.Headers["X-Custom"] != "value" {
			t.Errorf("expected custom header value, got %v", intent.Headers["X-Custom"])
		}
	})

	t.Run("location header when status code is 30x", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("X-Erised-Status-Code", "302")
		headers.Set("X-Erised-Location", "https://example.com")

		intent := BuildErisedIntent(headers, nil)

		if intent.Headers["Location"] != "https://example.com" {
			t.Errorf("expected Location header, got %v", intent.Headers["Location"])
		}
	})

	t.Run("location header ignored when status code is not 30x", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("X-Erised-Status-Code", "200")
		headers.Set("X-Erised-Location", "https://example.com")

		intent := BuildErisedIntent(headers, nil)

		if _, ok := intent.Headers["Location"]; ok {
			t.Errorf("expected Location header to be omitted")
		}
	})

	t.Run("response file overrides data and status code on success", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("X-Erised-Data", "fallback data")
		headers.Set("X-Erised-Response-File", "foo.txt")
		headers.Set("X-Erised-Status-Code", "500")

		resolver := func(filename string) ([]byte, error) {
			if filename == "foo.txt" {
				return []byte("file content"), nil
			}
			return nil, ErrFileNotFound
		}

		intent := BuildErisedIntent(headers, resolver)

		if intent.Data != "file content" {
			t.Errorf("expected file content, got %v", intent.Data)
		}
		if intent.StatusCode != 200 {
			t.Errorf("expected status code 200, got %v", intent.StatusCode)
		}
	})

	t.Run("response file updates status code on error", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("X-Erised-Data", "fallback data")
		headers.Set("X-Erised-Response-File", "missing.txt")
		headers.Set("X-Erised-Status-Code", "200")

		resolver := func(filename string) ([]byte, error) {
			return nil, ErrFileNotFound
		}

		intent := BuildErisedIntent(headers, resolver)

		if intent.Data != "" {
			t.Errorf("expected empty data, got %v", intent.Data)
		}
		if intent.StatusCode != 404 {
			t.Errorf("expected status code 404, got %v", intent.StatusCode)
		}
	})
}
