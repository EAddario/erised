package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/franela/goblin"
)

func TestErisedInfoRoute(t *testing.T) {
	g := goblin.Goblin(t)
	exp := `{"Host":"localhost:8080","Method":"GET","Protocol":"HTTP/1.1","Request URI":"http://localhost:8080/erised/info"}`
	svr := server {}
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/erised/info", nil)
	res := httptest.NewRecorder()
	svr.handleInfo().ServeHTTP(res, req)

	g.Describe("Test erised/info", func() {
		g.It("Should return StatusOK", func() {
			g.Assert(res.Code).Equal(http.StatusOK)
		})

		g.It("Should match expected body", func() {
			g.Assert(res.Body.String()).Equal(exp)
		})

		g.It("Should match Content-Type header", func() {
			g.Assert(res.Header().Get("Content-Type")).Equal("application/json")
		})

	})
}

func TestErisedIPRoute(t *testing.T) {
	g := goblin.Goblin(t)
	exp := `{"Client IP":"192.0.2.1:1234"}`
	svr := server {}
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/erised/ip", nil)
	res := httptest.NewRecorder()
	svr.handleIP().ServeHTTP(res, req)

	g.Describe("Test erised/ip", func() {
		g.It("Should return StatusOK", func() {
			g.Assert(res.Code).Equal(http.StatusOK)
		})

		g.It("Should match expected body", func() {
			g.Assert(res.Body.String()).Equal(exp)
		})

		g.It("Should match Content-Type header", func() {
			g.Assert(res.Header().Get("Content-Type")).Equal("application/json")
		})
	})
}

func TestErisedHeadersRoute(t *testing.T) {
	g := goblin.Goblin(t)
	exp := `{"Host":"localhost:8080"}`
	svr := server {}
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/erised/headers", nil)
	res := httptest.NewRecorder()
	svr.handleHeaders().ServeHTTP(res, req)

	g.Describe("Test erised/headers", func() {
		g.It("Should return StatusOK", func() {
			g.Assert(res.Code).Equal(http.StatusOK)
		})

		g.It("Should match expected body", func() {
			g.Assert(res.Body.String()).Equal(exp)
		})

		g.It("Should match Content-Type header", func() {
			g.Assert(res.Header().Get("Content-Type")).Equal("application/json")
		})
	})
}

func TestErisedLandingRoute(t *testing.T) {
	g := goblin.Goblin(t)
	svr := server {}

	g.Describe("Test /", func() {
		g.It("Should return StatusOK", func() {
			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
			svr.handleLanding().ServeHTTP(res, req)

			g.Assert(res.Code).Equal(http.StatusOK)
			g.Assert(res.Body.String()).Equal("")
		})

		g.It("Should return TemporaryRedirect and Location", func() {
			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
			req.Header.Set("X-Erised-Status-Code", "TemporaryRedirect")
			req.Header.Set("X-Erised-Location", "https://www.example.com")
			svr.handleLanding().ServeHTTP(res, req)

			g.Assert(res.Code).Equal(http.StatusTemporaryRedirect)
			g.Assert(res.Header().Get("Location")).Equal("https://www.example.com")
		})

		g.It("Should return JSON body", func() {
			exp := `{"hello":"world"}`
			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
			req.Header.Set("X-Erised-Content-Type", "json")
			req.Header.Set("X-Erised-Data", exp)
			svr.handleLanding().ServeHTTP(res, req)

			g.Assert(res.Code).Equal(http.StatusOK)
			g.Assert(res.Header().Get("Content-Type")).Equal("application/json")
			g.Assert(res.Body.String()).Equal(exp)
		})

		g.It("Should return text body", func() {
			exp := "Lorem ipsum"
			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
			req.Header.Set("X-Erised-Content-Type", "text")
			req.Header.Set("X-Erised-Data", exp)
			svr.handleLanding().ServeHTTP(res, req)

			g.Assert(res.Code).Equal(http.StatusOK)
			g.Assert(res.Header().Get("Content-Type")).Equal("text/plain")
			g.Assert(res.Body.String()).Equal(exp)
		})

		g.It("Should return headers", func() {
			exp := `{"X-Headers-One":"I'm header one","X-Headers-Two":"I'm header two"}`
			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
			req.Header.Set("X-Erised-Headers", exp)
			svr.handleLanding().ServeHTTP(res, req)

			g.Assert(res.Code).Equal(http.StatusOK)
			g.Assert(res.Header().Get("X-Headers-One")).Equal("I'm header one")
			g.Assert(res.Header().Get("X-Headers-Two")).Equal("I'm header two")
		})

		g.It("Should wait 2000ms", func() {
			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
			req.Header.Set("X-Erised-Response-Delay", "2000")
			svr.handleLanding().ServeHTTP(res, req)

			g.Assert(res.Code).Equal(http.StatusOK)
			g.Assert(res.Body.String()).Equal("")
		})

		g.It("Should not fail", func() {
			exp := `{"hello":"world"}`
			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
			req.Header.Set("X-Erised-Content-Type", "json")
			req.Header.Set("X-Erised-Data", exp)
			req.Header.Set("X-Erised-Headers", exp)
			req.Header.Set("X-Erised-Location", "https://www.example.com")
			req.Header.Set("X-Erised-Response-Delay", "1")
			req.Header.Set("X-Erised-Status-Code", "MovedPermanently")
			svr.handleLanding().ServeHTTP(res, req)

			g.Assert(res.Code).Equal(http.StatusMovedPermanently)
			g.Assert(res.Header().Get("Location")).Equal("https://www.example.com")
			g.Assert(res.Header().Get("Content-Type")).Equal("application/json")
			g.Assert(res.Header().Get("Content-Encoding")).Equal("identity")
			g.Assert(res.Header().Get("hello")).Equal("world")
			g.Assert(res.Body.String()).Equal(exp)
		})

	})
}

