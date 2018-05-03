package httputil

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/phogolabs/rho/httperr"
	uuid "github.com/satori/go.uuid"
)

// URLQueryParam returns the query param
func URLQueryParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

// URLQueryParamUUID returns a request query parameter as UUID
func URLQueryParamUUID(r *http.Request, key string) (uuid.UUID, error) {
	param := URLQueryParam(r, key)

	if param == "" {
		return uuid.Nil, httperr.QueryParamRequired(key)
	}

	value, err := uuid.FromString(param)
	if err == nil {
		return value, nil
	}

	err = httperr.QueryParamParse(key, "UUID", err)
	return uuid.Nil, err
}

// URLQueryParamUUIDOrValue returns a request query parameter as UUID or the
// provided default value if cannot parse the parameter.
func URLQueryParamUUIDOrValue(r *http.Request, key string, value uuid.UUID) uuid.UUID {
	param, err := URLQueryParamUUID(r, key)
	if err != nil {
		param = value
	}

	return param
}

// URLQueryParamUUIDOrNil returns a nil value if cannot parse the UUID parameter
func URLQueryParamUUIDOrNil(r *http.Request, key string) uuid.UUID {
	return URLQueryParamUUIDOrValue(r, key, uuid.Nil)
}

// URLQueryParamInt returns a request query parameter as int64
func URLQueryParamInt(r *http.Request, key string, base, bitSize int) (int64, error) {
	param := URLQueryParam(r, key)

	if param == "" {
		return 0, httperr.QueryParamRequired(key)
	}

	value, err := strconv.ParseInt(param, base, bitSize)
	if err == nil {
		return value, nil
	}

	err = httperr.QueryParamParse(key, "integer number", err)
	return 0, err
}

// URLQueryParamIntOrValue returns a request query parameter as int64 or the
// provided default value if cannot parse the parameter.
func URLQueryParamIntOrValue(r *http.Request, key string, base, bitSize int, value int64) int64 {
	param, err := URLQueryParamInt(r, key, base, bitSize)
	if err != nil {
		param = value
	}

	return param
}

// URLQueryParamUint returns a request query parameter as uint64
func URLQueryParamUint(r *http.Request, key string, base, bitSize int) (uint64, error) {
	param := URLQueryParam(r, key)

	if param == "" {
		return 0, httperr.QueryParamRequired(key)
	}

	value, err := strconv.ParseUint(param, base, bitSize)
	if err == nil {
		return value, nil
	}

	err = httperr.QueryParamParse(key, "unsigned integer number", err)
	return 0, err
}

// URLQueryParamUintOrValue returns a request query parameter as uint64 or the
// provided default value if cannot parse the parameter.
func URLQueryParamUintOrValue(r *http.Request, key string, base, bitSize int, value uint64) uint64 {
	param, err := URLQueryParamUint(r, key, base, bitSize)
	if err != nil {
		param = value
	}

	return param
}

// URLQueryParamFloat returns a request query parameter as float64
func URLQueryParamFloat(r *http.Request, key string, bitSize int) (float64, error) {
	param := URLQueryParam(r, key)

	if param == "" {
		return 0, httperr.QueryParamRequired(key)
	}

	value, err := strconv.ParseFloat(param, bitSize)
	if err == nil {
		return value, nil
	}

	err = httperr.QueryParamParse(key, "float number", err)
	return 0, err
}

// URLQueryParamFloatOrValue returns a request query parameter as float64 or the
// provided default value if cannot parse the parameter.
func URLQueryParamFloatOrValue(r *http.Request, key string, bitSize int, value float64) float64 {
	param, err := URLQueryParamFloat(r, key, bitSize)
	if err != nil {
		param = value
	}

	return param
}

// URLQueryParamTime returns a request query parameter as time.Time
func URLQueryParamTime(r *http.Request, key, format string) (time.Time, error) {
	param := URLQueryParam(r, key)

	if param == "" {
		return time.Time{}, httperr.QueryParamRequired(key)
	}

	value, err := time.Parse(format, param)
	if err == nil {
		return value, nil
	}

	info := fmt.Sprintf("Expected date time format '%s'", format)
	err = httperr.QueryParamParse(key, "date time", err, info)
	return time.Time{}, err
}

// URLQueryParamTimeOrValue returns a request query parameter as time.Time or the
// provided default value if cannot parse the parameter.
func URLQueryParamTimeOrValue(r *http.Request, key, format string, value time.Time) time.Time {
	param, err := URLQueryParamTime(r, key, format)
	if err != nil {
		param = value
	}

	return param
}