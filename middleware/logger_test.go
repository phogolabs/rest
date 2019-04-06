package middleware_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/go-chi/chi"
	"github.com/onsi/gomega/gbytes"
	"github.com/phogolabs/log"
	"github.com/phogolabs/log/handler/json"
	"github.com/phogolabs/rest/middleware"
)

var _ = Describe("Logger", func() {
	var output *gbytes.Buffer

	BeforeEach(func() {
		output = gbytes.NewBuffer()
		log.SetHandler(json.New(output))
	})

	It("writes to the log", func() {
		router := chi.NewMux()
		router.Use(middleware.Logger)

		handler := func(w http.ResponseWriter, r *http.Request) {
			middleware.GetLogger(r).Info("hello")
		}

		router.Mount("/", http.HandlerFunc(handler))
		router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "http://example.com/", nil))

		Expect(output).To(gbytes.Say("hello"))
	})
})
