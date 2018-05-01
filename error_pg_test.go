package rho_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/lib/pq"
	"github.com/phogolabs/rho"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PostgreSQL Error", func() {
	var (
		r *http.Request
		w *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		w = httptest.NewRecorder()
		r = httptest.NewRequest("GET", "http://example.com", nil)
	})

	Context("Class 08 - Connection Exception", func() {
		It("handles the error correctly", func() {
			err := pq.Error{Code: "08P01"}
			rho.HandleErr(w, r, err)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(rho.ErrBackendNotConnected)))
			Expect(payload).To(HaveKeyWithValue("message", "Connection Error"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})

	Context("Class 22 - Data Exception", func() {
		It("handles the error correctly", func() {
			err := pq.Error{Code: "22001"}
			rho.HandleErr(w, r, err)

			Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(rho.ErrConflict)))
			Expect(payload).To(HaveKeyWithValue("message", "Data Error"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})

		Context("when the error is numeric_value_out_of_range", func() {
			It("handles the error correctly", func() {
				err := pq.Error{Code: "22003"}
				rho.HandleErr(w, r, err)

				Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))
				payload := unmarshalErrResponse(w.Body)

				Expect(payload).To(HaveKeyWithValue("code", float64(rho.ErrOutOfrange)))
				Expect(payload).To(HaveKeyWithValue("message", "Data Error"))
				Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
			})
		})

		Context("when the error is datetime_field_overflow", func() {
			It("handles the error correctly", func() {
				err := pq.Error{Code: "22008"}
				rho.HandleErr(w, r, err)

				Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))
				payload := unmarshalErrResponse(w.Body)

				Expect(payload).To(HaveKeyWithValue("code", float64(rho.ErrOutOfrange)))
				Expect(payload).To(HaveKeyWithValue("message", "Data Error"))
				Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
			})
		})

		Context("when the error is interval_field_overflow", func() {
			It("handles the error correctly", func() {
				err := pq.Error{Code: "22015"}
				rho.HandleErr(w, r, err)

				Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))
				payload := unmarshalErrResponse(w.Body)

				Expect(payload).To(HaveKeyWithValue("code", float64(rho.ErrOutOfrange)))
				Expect(payload).To(HaveKeyWithValue("message", "Data Error"))
				Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
			})
		})

		Context("when the error is indicator_overflow", func() {
			It("handles the error correctly", func() {
				err := pq.Error{Code: "22022"}
				rho.HandleErr(w, r, err)

				Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))
				payload := unmarshalErrResponse(w.Body)

				Expect(payload).To(HaveKeyWithValue("code", float64(rho.ErrOutOfrange)))
				Expect(payload).To(HaveKeyWithValue("message", "Data Error"))
				Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
			})
		})

		Context("when the error is floating_point_exception", func() {
			It("handles the error correctly", func() {
				err := pq.Error{Code: "22P01"}
				rho.HandleErr(w, r, err)

				Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))
				payload := unmarshalErrResponse(w.Body)

				Expect(payload).To(HaveKeyWithValue("code", float64(rho.ErrOutOfrange)))
				Expect(payload).To(HaveKeyWithValue("message", "Data Error"))
				Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
			})
		})

		Context("when the error is null_value_not_allowed", func() {
			It("handles the error correctly", func() {
				err := pq.Error{Code: "22002"}
				rho.HandleErr(w, r, err)

				Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))
				payload := unmarshalErrResponse(w.Body)

				Expect(payload).To(HaveKeyWithValue("code", float64(rho.ErrConditionNotMet)))
				Expect(payload).To(HaveKeyWithValue("message", "Data Error"))
				Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
			})
		})

		Context("when the error is null_value_no_indicator_parameter", func() {
			It("handles the error correctly", func() {
				err := pq.Error{Code: "22004"}
				rho.HandleErr(w, r, err)

				Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))
				payload := unmarshalErrResponse(w.Body)

				Expect(payload).To(HaveKeyWithValue("code", float64(rho.ErrConditionNotMet)))
				Expect(payload).To(HaveKeyWithValue("message", "Data Error"))
				Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
			})
		})
	})

	Context("Class 23 - Integrity Constraint Violation", func() {
		It("handles the error correctly", func() {
		})

		Context("when the error is unique_violation", func() {
			It("handles the error correctly", func() {
				err := pq.Error{Code: "23505"}
				rho.HandleErr(w, r, err)

				Expect(w.Code).To(Equal(http.StatusConflict))
				payload := unmarshalErrResponse(w.Body)

				Expect(payload).To(HaveKeyWithValue("code", float64(rho.ErrDuplicate)))
				Expect(payload).To(HaveKeyWithValue("message", "Integrity Constraint Violation"))
				Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
			})
		})

		Context("when the error is check_violation", func() {
			It("handles the error correctly", func() {
				err := pq.Error{Code: "23514"}
				rho.HandleErr(w, r, err)

				Expect(w.Code).To(Equal(http.StatusConflict))
				payload := unmarshalErrResponse(w.Body)

				Expect(payload).To(HaveKeyWithValue("code", float64(rho.ErrConditionNotMet)))
				Expect(payload).To(HaveKeyWithValue("message", "Integrity Constraint Violation"))
				Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
			})
		})

		Context("when the error is exclusion_violation", func() {
			It("handles the error correctly", func() {
				err := pq.Error{Code: "23P01"}
				rho.HandleErr(w, r, err)

				Expect(w.Code).To(Equal(http.StatusConflict))
				payload := unmarshalErrResponse(w.Body)

				Expect(payload).To(HaveKeyWithValue("code", float64(rho.ErrConditionNotMet)))
				Expect(payload).To(HaveKeyWithValue("message", "Integrity Constraint Violation"))
				Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
			})
		})
	})

	Context("Class 57 - Operation Intervation", func() {
		It("handles the error correctly", func() {
			err := pq.Error{Code: "57P01"}
			rho.HandleErr(w, r, err)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(rho.ErrNotReady)))
			Expect(payload).To(HaveKeyWithValue("message", "Operator Intervention"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})

	Context("When the Class is unknown", func() {
		It("handles the error correctly", func() {
			err := pq.Error{Code: "9999"}
			rho.HandleErr(w, r, err)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(rho.ErrBackend)))
			Expect(payload).To(HaveKeyWithValue("message", "Database Error"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})
})
