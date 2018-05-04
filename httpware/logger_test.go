package httpware_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/rho/httpware"
)

var _ = Describe("HTTPWare", func() {
	var (
		r      *http.Request
		w      *httptest.ResponseRecorder
		h      http.Handler
		buffer *bytes.Buffer
		code   int
	)

	BeforeEach(func() {
		code = http.StatusOK

		buffer = &bytes.Buffer{}
		log.SetHandler(text.New(buffer))
	})

	JustBeforeEach(func() {
		w = httptest.NewRecorder()
		r = httptest.NewRequest("GET", "http://example.com", nil)

		h = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(code)
			w.Write([]byte("hello"))
		})

		h = httpware.Logger(h)
	})

	It("logs the request on info level", func() {
		h.ServeHTTP(w, r)
		Expect(w.Code).To(Equal(http.StatusOK))
		Expect(w.Body.String()).To(Equal("hello"))
		Expect(buffer.String()).To(ContainSubstring("INFO"))
	})

	Context("when the code is greater than 400", func() {
		BeforeEach(func() {
			code = http.StatusBadRequest
		})

		It("logs the request on warn level", func() {
			h.ServeHTTP(w, r)
			Expect(w.Code).To(Equal(http.StatusBadRequest))
			Expect(w.Body.String()).To(Equal("hello"))
			Expect(buffer.String()).To(ContainSubstring("WARN"))
		})
	})

	Context("when the code is greater than 500", func() {
		BeforeEach(func() {
			code = http.StatusInternalServerError
		})

		It("logs the request on error level", func() {
			h.ServeHTTP(w, r)
			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			Expect(w.Body.String()).To(Equal("hello"))
			Expect(buffer.String()).To(ContainSubstring("ERROR"))
		})
	})

	Context("when the panic is logged", func() {
		It("logs the request on error level", func() {
			entry := httpware.NewLogEntry(r)
			entry.Panic(fmt.Errorf("Oh no!"), nil)
			Expect(buffer.String()).To(ContainSubstring("Oh no!"))
			Expect(buffer.String()).To(ContainSubstring("ERROR"))
		})

		Context("when the panic is not an error", func() {
			It("logs the request on error level", func() {
				entry := httpware.NewLogEntry(r)
				entry.Panic(-1, []byte("lol"))
				Expect(buffer.String()).To(ContainSubstring("-1"))
				Expect(buffer.String()).To(ContainSubstring("lol"))
			})
		})
	})
})

var _ = Describe("SetLogger", func() {
	var logger log.Interface

	BeforeEach(func() {
		logger = log.Log
	})

	AfterEach(func() {
		log.Log = logger
	})

	ItConfiguresTheLoggerSuccessfully := func(format string) {
		Context(fmt.Sprintf("when the format is '%s'", format), func() {
			It("configures the logger successfully", func() {
				buffer := &bytes.Buffer{}

				cfg := &httpware.LoggerConfig{
					Fields: log.Fields{
						"app_name":    "test",
						"app_version": "beta",
					},
					Format: format,
					Level:  "info",
					Output: buffer,
				}

				Expect(httpware.SetLogger(cfg)).To(Succeed())

				log.Info("hello")

				Expect(buffer.String()).To(ContainSubstring("test"))
				Expect(buffer.String()).To(ContainSubstring("beta"))
			})
		})
	}

	ItConfiguresTheLoggerSuccessfully("json")
	ItConfiguresTheLoggerSuccessfully("text")
	ItConfiguresTheLoggerSuccessfully("cli")

	Context("when the format is not supported", func() {
		It("returns an error", func() {
			buffer := &bytes.Buffer{}

			cfg := &httpware.LoggerConfig{
				Format: "wrong",
				Output: buffer,
			}

			Expect(httpware.SetLogger(cfg)).To(MatchError("unsupported log format 'wrong'"))
		})
	})

	Context("when the level is wrong", func() {
		It("returns an error", func() {
			buffer := &bytes.Buffer{}

			cfg := &httpware.LoggerConfig{
				Fields: log.Fields{
					"app_name":    "test",
					"app_version": "beta",
				},
				Format: "text",
				Level:  "wrong",
				Output: buffer,
			}

			Expect(httpware.SetLogger(cfg)).To(MatchError("invalid level"))
		})
	})
})
