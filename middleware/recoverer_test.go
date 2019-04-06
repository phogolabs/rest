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

var _ = Describe("Recoverer", func() {
	var output *gbytes.Buffer

	BeforeEach(func() {
		output = gbytes.NewBuffer()
		log.SetHandler(json.New(output))
	})

	It("recovers sucessfully", func() {
		router := chi.NewMux()
		router.Use(middleware.Recoverer)

		handler := func(w http.ResponseWriter, r *http.Request) {
			panic("hello")
		}

		router.Mount("/", http.HandlerFunc(handler))
		router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "http://example.com/", nil))

		Expect(output).To(gbytes.Say("hello"))
	})
})
