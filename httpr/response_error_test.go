package httpr_test

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/lib/pq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/http/httpr"
	validator "gopkg.in/go-playground/validator.v9"
)

var _ = Describe("Conv Error", func() {
	var (
		r *http.Request
		w *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		w = httptest.NewRecorder()
		r = httptest.NewRequest("GET", "http://example.com", nil)
	})

	Context("when a time error occurs", func() {
		It("handles the error correctly", func() {
			err := &time.ParseError{}
			httpr.RespondError(w, r, err)

			Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(httpr.CodeInvalid)))
			Expect(payload).To(HaveKeyWithValue("message", "Unable to parse date time"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})

	Context("when a num error occurs", func() {
		It("handles the error correctly", func() {
			err := &strconv.NumError{Err: fmt.Errorf("oh no!")}
			httpr.RespondError(w, r, err)

			Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(httpr.CodeInvalid)))
			Expect(payload).To(HaveKeyWithValue("message", "Unable to parse number"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})
})

var _ = Describe("JSON Error", func() {
	var (
		r *http.Request
		w *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		w = httptest.NewRecorder()
		r = httptest.NewRequest("GET", "http://example.com", nil)
	})

	Context("when the err is not recognized", func() {
		It("handles the error correctly", func() {
			errx := httpr.JSONError(fmt.Errorf("Oh no!"))
			Expect(errx.StatusCode).To(Equal(http.StatusInternalServerError))

			herr, ok := errx.Err.(*httpr.HTTPError)
			Expect(ok).To(BeTrue())
			Expect(herr.Code).To(Equal(httpr.CodeInternal))
			Expect(herr.Message).To(Equal("JSON Error"))
		})
	})

	Context("when the err is json.InvalidUnmarshalError", func() {
		It("handles the error correctly", func() {
			err := &json.InvalidUnmarshalError{Type: reflect.TypeOf(r)}
			httpr.RespondError(w, r, err)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(httpr.CodeInternal)))
			Expect(payload).To(HaveKeyWithValue("message", "Unable to unmarshal json body"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})

	Context("when the err is json.UnmarshalFieldError", func() {
		It("handles the error correctly", func() {
			typ := reflect.TypeOf(*r)
			err := &json.UnmarshalFieldError{Type: typ, Field: typ.Field(0), Key: "StatusCode"}
			httpr.RespondError(w, r, err)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(httpr.CodeInvalid)))
			Expect(payload).To(HaveKeyWithValue("message", "Unable to unmarshal json body"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})

	Context("when the err is UnmarshalTypeError", func() {
		It("handles the error correctly", func() {
			err := &json.UnmarshalTypeError{Type: reflect.TypeOf(*r)}
			httpr.RespondError(w, r, err)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(httpr.CodeInvalid)))
			Expect(payload).To(HaveKeyWithValue("message", "Unable to unmarshal json body"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})

	Context("when the err is UnsupportedValueError", func() {
		It("handles the error correctly", func() {
			err := &json.UnsupportedValueError{Value: reflect.ValueOf(*r)}
			httpr.RespondError(w, r, err)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(httpr.CodeInternal)))
			Expect(payload).To(HaveKeyWithValue("message", "Unable to marshal json"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})

	Context("when the err is UnsupportedTypeError", func() {
		It("handles the error correctly", func() {
			err := &json.UnsupportedTypeError{Type: reflect.TypeOf(*r)}
			httpr.RespondError(w, r, err)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(httpr.CodeInternal)))
			Expect(payload).To(HaveKeyWithValue("message", "Unable to marshal json"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})

	Context("when the err is InvalidUTF8Error", func() {
		It("handles the error correctly", func() {
			err := &json.InvalidUTF8Error{S: "Oh no!"}
			httpr.RespondError(w, r, err)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(httpr.CodeInternal)))
			Expect(payload).To(HaveKeyWithValue("message", "Unable to marshal json"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})

	Context("when the err is MarshalerError", func() {
		It("handles the error correctly", func() {
			err := &json.MarshalerError{Err: fmt.Errorf("Oh no!"), Type: reflect.TypeOf(*r)}
			httpr.RespondError(w, r, err)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(httpr.CodeInternal)))
			Expect(payload).To(HaveKeyWithValue("message", "Unable to marshal json"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})
})

var _ = Describe("XML Error", func() {
	var (
		r *http.Request
		w *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		w = httptest.NewRecorder()
		r = httptest.NewRequest("GET", "http://example.com", nil)
	})

	Context("when the err is not recognized", func() {
		It("handles the error correctly", func() {
			errx := httpr.XMLError(fmt.Errorf("Oh no!"))
			herr, ok := errx.Err.(*httpr.HTTPError)
			Expect(ok).To(BeTrue())
			Expect(errx.StatusCode).To(Equal(http.StatusInternalServerError))
			Expect(herr.Code).To(Equal(httpr.CodeInternal))
			Expect(herr.Message).To(Equal("XML Error"))
		})
	})

	Context("when the err is xml.UnmarshalError", func() {
		It("handles the error correctly", func() {
			err := xml.UnmarshalError("oh no")
			httpr.RespondError(w, r, err)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(httpr.CodeInvalid)))
			Expect(payload).To(HaveKeyWithValue("message", "Unable to unmarshal xml body"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})

	Context("when the err is xml.SyntaxError", func() {
		It("handles the error correctly", func() {
			err := &xml.SyntaxError{Msg: "oh no!"}
			httpr.RespondError(w, r, err)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(httpr.CodeInvalid)))
			Expect(payload).To(HaveKeyWithValue("message", "Unable to unmarshal xml body"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})

	Context("when the err is xml.TagPathError", func() {
		It("handles the error correctly", func() {
			err := &xml.TagPathError{Struct: reflect.TypeOf(*r)}
			httpr.RespondError(w, r, err)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(httpr.CodeInvalid)))
			Expect(payload).To(HaveKeyWithValue("message", "Unable to unmarshal xml body"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})

	Context("when the err is xml.UnsupportedTypeError", func() {
		It("handles the error correctly", func() {
			err := &xml.UnsupportedTypeError{Type: reflect.TypeOf(*r)}
			httpr.RespondError(w, r, err)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(httpr.CodeInternal)))
			Expect(payload).To(HaveKeyWithValue("message", "Unable to marshal xml"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})
})

var _ = Describe("Validation Error", func() {
	var (
		r *http.Request
		w *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		w = httptest.NewRecorder()
		r = httptest.NewRequest("GET", "http://example.com", nil)
	})

	Context("when the error is validation error", func() {
		It("respond with the provided information", func() {
			type User struct {
				Name     string `json:"name" validate:"required"`
				Password string `json:"password" validate:"required"`
			}

			v := validator.New()
			err := v.Struct(&User{})

			httpr.RespondError(w, r, err)

			Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))
			Expect(w.Body.String()).To(ContainSubstring("Field 'Name' is not valid"))
			Expect(w.Body.String()).To(ContainSubstring("Field 'Password' is not valid"))
		})
	})

	Context("when the error is InvalidValidationError", func() {
		It("handles the error correctoy", func() {
			err := &validator.InvalidValidationError{Type: reflect.TypeOf(*r)}
			resp := httpr.ValidationError(err)
			Expect(resp.StatusCode).To(Equal(http.StatusUnprocessableEntity))
			Expect(resp).To(MatchError("Validation Error"))
		})
	})
})

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
			httpr.RespondError(w, r, err)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(httpr.CodeBackendNotConnected)))
			Expect(payload).To(HaveKeyWithValue("message", "Connection Error"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})

	Context("Class 22 - Data Exception", func() {
		It("handles the error correctly", func() {
			err := pq.Error{Code: "22001"}
			httpr.RespondError(w, r, err)

			Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(httpr.CodeConflict)))
			Expect(payload).To(HaveKeyWithValue("message", "Data Error"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})

		Context("when the error is numeric_value_out_of_range", func() {
			It("handles the error correctly", func() {
				err := pq.Error{Code: "22003"}
				httpr.RespondError(w, r, err)

				Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))
				payload := unmarshalErrResponse(w.Body)

				Expect(payload).To(HaveKeyWithValue("code", float64(httpr.CodeOutOfrange)))
				Expect(payload).To(HaveKeyWithValue("message", "Data Error"))
				Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
			})
		})

		Context("when the error is datetime_field_overflow", func() {
			It("handles the error correctly", func() {
				err := pq.Error{Code: "22008"}
				httpr.RespondError(w, r, err)

				Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))
				payload := unmarshalErrResponse(w.Body)

				Expect(payload).To(HaveKeyWithValue("code", float64(httpr.CodeOutOfrange)))
				Expect(payload).To(HaveKeyWithValue("message", "Data Error"))
				Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
			})
		})

		Context("when the error is interval_field_overflow", func() {
			It("handles the error correctly", func() {
				err := pq.Error{Code: "22015"}
				httpr.RespondError(w, r, err)

				Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))
				payload := unmarshalErrResponse(w.Body)

				Expect(payload).To(HaveKeyWithValue("code", float64(httpr.CodeOutOfrange)))
				Expect(payload).To(HaveKeyWithValue("message", "Data Error"))
				Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
			})
		})

		Context("when the error is indicator_overflow", func() {
			It("handles the error correctly", func() {
				err := pq.Error{Code: "22022"}
				httpr.RespondError(w, r, err)

				Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))
				payload := unmarshalErrResponse(w.Body)

				Expect(payload).To(HaveKeyWithValue("code", float64(httpr.CodeOutOfrange)))
				Expect(payload).To(HaveKeyWithValue("message", "Data Error"))
				Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
			})
		})

		Context("when the error is floating_point_exception", func() {
			It("handles the error correctly", func() {
				err := pq.Error{Code: "22P01"}
				httpr.RespondError(w, r, err)

				Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))
				payload := unmarshalErrResponse(w.Body)

				Expect(payload).To(HaveKeyWithValue("code", float64(httpr.CodeOutOfrange)))
				Expect(payload).To(HaveKeyWithValue("message", "Data Error"))
				Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
			})
		})

		Context("when the error is null_value_not_allowed", func() {
			It("handles the error correctly", func() {
				err := pq.Error{Code: "22002"}
				httpr.RespondError(w, r, err)

				Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))
				payload := unmarshalErrResponse(w.Body)

				Expect(payload).To(HaveKeyWithValue("code", float64(httpr.CodeConditionNotMet)))
				Expect(payload).To(HaveKeyWithValue("message", "Data Error"))
				Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
			})
		})

		Context("when the error is null_value_no_indicator_parameter", func() {
			It("handles the error correctly", func() {
				err := pq.Error{Code: "22004"}
				httpr.RespondError(w, r, err)

				Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))
				payload := unmarshalErrResponse(w.Body)

				Expect(payload).To(HaveKeyWithValue("code", float64(httpr.CodeConditionNotMet)))
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
				httpr.RespondError(w, r, err)

				Expect(w.Code).To(Equal(http.StatusConflict))
				payload := unmarshalErrResponse(w.Body)

				Expect(payload).To(HaveKeyWithValue("code", float64(httpr.CodeDuplicate)))
				Expect(payload).To(HaveKeyWithValue("message", "Integrity Constraint Violation"))
				Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
			})
		})

		Context("when the error is check_violation", func() {
			It("handles the error correctly", func() {
				err := pq.Error{Code: "23514"}
				httpr.RespondError(w, r, err)

				Expect(w.Code).To(Equal(http.StatusConflict))
				payload := unmarshalErrResponse(w.Body)

				Expect(payload).To(HaveKeyWithValue("code", float64(httpr.CodeConditionNotMet)))
				Expect(payload).To(HaveKeyWithValue("message", "Integrity Constraint Violation"))
				Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
			})
		})

		Context("when the error is exclusion_violation", func() {
			It("handles the error correctly", func() {
				err := pq.Error{Code: "23P01"}
				httpr.RespondError(w, r, err)

				Expect(w.Code).To(Equal(http.StatusConflict))
				payload := unmarshalErrResponse(w.Body)

				Expect(payload).To(HaveKeyWithValue("code", float64(httpr.CodeConditionNotMet)))
				Expect(payload).To(HaveKeyWithValue("message", "Integrity Constraint Violation"))
				Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
			})
		})
	})

	Context("Class 57 - Operation", func() {
		It("handles the error correctly", func() {
			err := pq.Error{Code: "57P01"}
			httpr.RespondError(w, r, err)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(httpr.CodeBackendNotReady)))
			Expect(payload).To(HaveKeyWithValue("message", "Operator Intervention"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})

	Context("When the Class is unknown", func() {
		It("handles the error correctly", func() {
			err := pq.Error{Code: "9999"}
			httpr.RespondError(w, r, err)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(httpr.CodeBackend)))
			Expect(payload).To(HaveKeyWithValue("message", "Database Error"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})
})

var _ = Describe("ErrorResponse", func() {
	var (
		r *http.Request
		w *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		w = httptest.NewRecorder()
		r = httptest.NewRequest("GET", "http://example.com", nil)
	})

	Context("when the error is http error", func() {
		It("handles the error", func() {
			err := httpr.NewError(1, "Oh no!")
			httpr.RespondError(w, r, err)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(1)))
			Expect(payload).To(HaveKeyWithValue("message", "Oh no!"))
		})
	})

	Context("when the error is nil", func() {
		It("handles the error", func() {
			httpr.RespondError(w, r, nil)

			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Body.Len()).To(BeZero())
		})
	})

	Context("when the error is regular error", func() {
		It("handles the error", func() {
			err := fmt.Errorf("Oh no!")
			httpr.RespondError(w, r, err)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(httpr.CodeInternal)))
			Expect(payload).To(HaveKeyWithValue("message", "Internal Error"))
			Expect(payload["reason"]).To(HaveKeyWithValue("message", err.Error()))
		})
	})

	Context("when the error is ErrorReponse", func() {
		It("handles the error", func() {
			err := &httpr.ErrorResponse{
				Err: httpr.NewError(1, "Oh no!"),
			}

			fn := func(ww http.ResponseWriter, rr *http.Request) {
				httpr.RespondError(ww, rr, err)
			}

			middleware.Logger(http.HandlerFunc(fn)).ServeHTTP(w, r)

			Expect(w.Code).To(Equal(err.StatusCode))
			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			payload := unmarshalErrResponse(w.Body)

			Expect(payload).To(HaveKeyWithValue("code", float64(1)))
			Expect(payload).To(HaveKeyWithValue("message", "Oh no!"))
			Expect(payload).NotTo(HaveKey("reason"))
		})
	})
})
