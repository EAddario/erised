package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog"
)

func TestErisedInfoRoute(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })
	exp := `{"Host":"localhost:8080","Method":"GET","Protocol":"HTTP/1.1","Request URI":"http://localhost:8080/erised/info"}`
	svr := server{}
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/erised/info", nil)
	res := httptest.NewRecorder()
	svr.handleInfo().ServeHTTP(res, req)

	g.Describe("Test erised/info", func() {
		g.It("Should return StatusOK", func() {
			Ω(res.Code).Should(Equal(http.StatusOK))
		})

		g.It("Should match expected body", func() {
			Ω(res.Body.String()).Should(Equal(exp))
		})

		g.It("Should match Content-Type header", func() {
			Ω(res.Header().Get("Content-Type")).Should(Equal("application/json"))
		})
	})
}

func TestErisedIPRoute(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })
	exp := `{"Client IP":"192.0.2.1:1234"}`
	svr := server{}
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/erised/ip", nil)
	res := httptest.NewRecorder()
	svr.handleIP().ServeHTTP(res, req)

	g.Describe("Test erised/ip", func() {
		g.It("Should return StatusOK", func() {
			Ω(res.Code).Should(Equal(http.StatusOK))
		})

		g.It("Should match expected body", func() {
			Ω(res.Body.String()).Should(Equal(exp))
		})

		g.It("Should match Content-Type header", func() {
			Ω(res.Header().Get("Content-Type")).Should(Equal("application/json"))
		})
	})
}

func TestErisedHeadersRoute(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })
	exp := `{"Host":"localhost:8080"}`
	svr := server{}
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/erised/headers", nil)
	res := httptest.NewRecorder()
	svr.handleHeaders().ServeHTTP(res, req)

	g.Describe("Test erised/headers", func() {
		g.It("Should return StatusOK", func() {
			Ω(res.Code).Should(Equal(http.StatusOK))
		})

		g.It("Should match expected body", func() {
			Ω(res.Body.String()).Should(Equal(exp))
		})

		g.It("Should match Content-Type header", func() {
			Ω(res.Header().Get("Content-Type")).Should(Equal("application/json"))
		})
	})
}

func TestErisedLandingRoute(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })
	path := "."
	svr := server{pth: &path}

	g.Describe("Test /", func() {
		g.It("Should return StatusOK", func() {
			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
			svr.handleLanding().ServeHTTP(res, req)

			Ω(res).Should(HaveHTTPStatus(http.StatusOK))
			Ω(res.Body.String()).Should(BeEmpty())
		})

		g.It("Should return TemporaryRedirect and Location url", func() {
			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
			req.Header.Set("X-Erised-Status-Code", "TemporaryRedirect")
			req.Header.Set("X-Erised-Location", "https://www.example.com")
			svr.handleLanding().ServeHTTP(res, req)

			Ω(res).Should(HaveHTTPStatus(http.StatusTemporaryRedirect))
			Ω(res.Header().Get("Location")).Should(Equal("https://www.example.com"))
		})

		g.It("Should return json body", func() {
			exp := `{"hello":"world"}`
			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
			req.Header.Set("X-Erised-Content-Type", "json")
			req.Header.Set("X-Erised-Data", exp)
			svr.handleLanding().ServeHTTP(res, req)

			Ω(res).Should(HaveHTTPStatus(http.StatusOK))
			Ω(res.Header().Get("Content-Type")).Should(Equal("application/json"))
			Ω(res.Header().Get("Content-Encoding")).Should(Equal("identity"))
			Ω(res.Body.String()).Should(Equal(exp))
		})

		g.It("Should return text body", func() {
			exp := "Lorem ipsum"
			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
			req.Header.Set("X-Erised-Content-Type", "text")
			req.Header.Set("X-Erised-Data", exp)
			svr.handleLanding().ServeHTTP(res, req)

			Ω(res).Should(HaveHTTPStatus(http.StatusOK))
			Ω(res.Header().Get("Content-Type")).Should(Equal("text/plain"))
			Ω(res.Header().Get("Content-Encoding")).Should(Equal("identity"))
			Ω(res.Body.String()).Should(Equal(exp))
		})

		g.It("Should return xml body", func() {
			exp := "<hello>world</hello>"
			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
			req.Header.Set("X-Erised-Content-Type", "xml")
			req.Header.Set("X-Erised-Data", exp)
			svr.handleLanding().ServeHTTP(res, req)

			Ω(res).Should(HaveHTTPStatus(http.StatusOK))
			Ω(res.Header().Get("Content-Type")).Should(Equal("application/xml"))
			Ω(res.Header().Get("Content-Encoding")).Should(Equal("identity"))
			Ω(res.Body.String()).Should(Equal(exp))
		})

		g.It("Should return gzip body", func() {
			exp := "Lorem ipsum"
			gz := "\x1f\x8b\b\x00\x00\x00\x00\x00\x00\xff\xf2\xc9/J\xcdU\xc8,(.\xcd\x05\x04\x00\x00\xff\xffY\xfbK\xf4\v\x00\x00\x00"

			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
			req.Header.Set("X-Erised-Content-Type", "gzip")
			req.Header.Set("X-Erised-Data", exp)
			svr.handleLanding().ServeHTTP(res, req)

			Ω(res).Should(HaveHTTPStatus(http.StatusOK))
			Ω(res.Header().Get("Content-Type")).Should(Equal("application/octet-stream"))
			Ω(res.Header().Get("Content-Encoding")).Should(Equal("gzip"))
			Ω(res.Body.String()).Should(Equal(gz))
		})

		g.It("Should return headers", func() {
			exp := `{"X-Headers-One":"I'm header one","X-Headers-Two":"I'm header two"}`
			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
			req.Header.Set("X-Erised-Headers", exp)
			svr.handleLanding().ServeHTTP(res, req)

			Ω(res).Should(HaveHTTPStatus(http.StatusOK))
			Ω(res.Header().Get("X-Headers-One")).Should(Equal("I'm header one"))
			Ω(res.Header().Get("X-Headers-Two")).Should(Equal("I'm header two"))
		})

		g.It("Should return serverRoutes_test.json file content in body", func() {
			exp := `{"Name":"serverRoutes_test"}`
			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
			req.Header.Set("X-Erised-Content-Type", "json")
			req.Header.Set("X-Erised-Response-File", "serverRoutes_test.json")
			svr.handleLanding().ServeHTTP(res, req)

			Ω(res).Should(HaveHTTPStatus(http.StatusOK))
			Ω(res.Header().Get("Content-Type")).Should(Equal("application/json"))
			Ω(res.Header().Get("Content-Encoding")).Should(Equal("identity"))
			Ω(res.Body.String()).Should(Equal(exp))
		})

		g.It("Should not fail", func() {
			exp := `{"hello":"world"}`
			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
			req.Header.Set("X-Erised-Content-Type", "json")
			req.Header.Set("X-Erised-Data", exp)
			req.Header.Set("X-Erised-Headers", exp)
			req.Header.Set("X-Erised-Location", "https://www.example.com")
			req.Header.Set("X-Erised-Status-Code", "MovedPermanently")
			svr.handleLanding().ServeHTTP(res, req)

			Ω(res).Should(HaveHTTPStatus(http.StatusMovedPermanently))
			Ω(res.Header().Get("Location")).Should(Equal("https://www.example.com"))
			Ω(res.Header().Get("Content-Type")).Should(Equal("application/json"))
			Ω(res.Header().Get("Content-Encoding")).Should(Equal("identity"))
			Ω(res.Header().Get("hello")).Should(Equal("world"))
			Ω(res.Body.String()).Should(Equal(exp))
		})

		g.It("Should wait about 2000ms (±10ms)", func() {
			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
			req.Header.Set("X-Erised-Response-Delay", "2000")

			st := time.Now()
			svr.handleLanding().ServeHTTP(res, req)
			el := time.Since(st)

			Ω(res).Should(HaveHTTPStatus(http.StatusOK))
			Ω(el).Should(BeNumerically("~", time.Millisecond*2000, time.Millisecond*10))
			Ω(res.Body.String()).Should(BeEmpty())
		})
	})
}
