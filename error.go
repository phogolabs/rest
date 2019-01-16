package rest

import (
	"net/http"

	"github.com/apex/log"
	"github.com/go-chi/render"
	"github.com/go-playground/errors"
	"github.com/goware/errorx"
	multierror "github.com/hashicorp/go-multierror"
	rollbar "github.com/rollbar/rollbar-go"
	validator "gopkg.in/go-playground/validator.v9"
)

// Error injects the error within the request
func Error(w http.ResponseWriter, r *http.Request, err error) {
	err = errorChain(r, err)

	errorReport(r, err)
	errorStatus(r, err)

	Respond(w, r, errorWrap(err))
}

// ErrorXML injects the error within the request
func ErrorXML(w http.ResponseWriter, r *http.Request, err error) {
	err = errorChain(r, err)

	errorReport(r, err)
	errorStatus(r, err)

	XML(w, r, errorWrap(err))
}

// ErrorJSON injects the error within the request
func ErrorJSON(w http.ResponseWriter, r *http.Request, err error) {
	err = errorChain(r, err)

	errorReport(r, err)
	errorStatus(r, err)

	JSON(w, r, errorWrap(err))
}

func errorChain(r *http.Request, err error) error {
	ch, ok := err.(errors.Chain)

	if !ok {
		ch = errors.Wrap(err, "request")
	}

	if code, ok := errors.LookupTag(ch, "status").(int); !ok {
		if code, ok = r.Context().Value(render.StatusCtxKey).(int); !ok {
			code = http.StatusInternalServerError
		}

		ch = ch.AddTag("status", code)
	}

	return ch
}

func errorReport(r *http.Request, err error) {
	status := errors.LookupTag(err, "status").(int)

	fields := log.Fields{
		"status": status,
	}

	logger := GetLogger(r).
		WithError(err).
		WithFields(fields)

	switch {
	case status >= 500:
		logger.Error("occurred")
	case status >= 400:
		logger.Warn("occurred")
	default:
		logger.Info("occurred")
	}

	if rollbar.Token() == "" {
		return
	}

	switch {
	case status >= 500:
		rollbar.RequestErrorWithExtras(rollbar.ERR, r, err, fields)
	case status >= 400:
		rollbar.RequestErrorWithExtras(rollbar.WARN, r, err, fields)
	}
}

func errorStatus(r *http.Request, err error) {
	status := errors.LookupTag(err, "status").(int)
	Status(r, status)
}

func errorWrap(err error) *errorx.Errorx {
	code := errors.LookupTag(err, "status").(int)
	errx := errorx.New(code, http.StatusText(code))
	err = errors.Cause(err)

	if errs, ok := err.(*multierror.Error); ok {
		for _, err := range errs.Errors {
			errx.Details = append(errx.Details, err.Error())
		}
	} else if verrs, ok := err.(validator.ValidationErrors); ok {
		for _, verr := range verrs {
			if ferr, ok := verr.(error); ok {
				errx.Details = append(errx.Details, ferr.Error())
			}
		}
	} else {
		errx.Details = append(errx.Details, err.Error())
	}
	return errx
}
