package rest

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/errors"
	"github.com/go-playground/validator/v10"
	"github.com/goware/errorx"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/phogolabs/log"
)

func errorf(r *http.Request, err error) error {
	err = errorChain(r, err)

	errorReport(r, err)
	errorStatus(r, err)

	return errorWrap(err)
}

func errorChain(r *http.Request, err error) error {
	ch, ok := err.(errors.Chain)

	if !ok {
		ch = errors.WrapSkipFrames(err, "request", 4)
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

	fields := log.Map{
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
