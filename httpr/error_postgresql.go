package httpr

import (
	"net/http"

	"github.com/lib/pq"
)

const (
	// Class 08 - Connection Exception
	pgConnClassErr = "08"
	// Class 22 - Data Exception
	pgDataClassErr = "22"
	// Class 23 - Integrity Constraint Violation
	pgContraintClassErr = "23"
	// Class 57 - Operator Intervention
	pgOpIntClassErr = "57"
)

// PGError creates a Response for given PostgreSQL error
func PGError(err error) *ErrorResponse {
	var (
		pgErr    = err.(pq.Error)
		response *ErrorResponse
	)

	switch pgErr.Code[:2] {
	case pgConnClassErr:
		response = &ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        NewError(CodeBackendNotConnected, "Connection Error").Wrap(err),
		}
	case pgDataClassErr:
		response = PGDataError(pgErr)
	case pgContraintClassErr:
		response = PGIntegrityError(pgErr)
	case pgOpIntClassErr:
		response = &ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        NewError(CodeBackendNotReady, "Operator Intervention").Wrap(err),
		}
	default:
		response = &ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        NewError(CodeBackend, "Database Error").Wrap(err),
		}
	}

	return response
}

// PGIntegrityError handles PG integrity errors
func PGIntegrityError(err pq.Error) *ErrorResponse {
	errx := NewError(CodeConflict, "Integrity Constraint Violation")

	switch err.Code {
	// "23505": "unique_violation",
	case "23505":
		errx.Code = CodeDuplicate
	// "23514": "check_violation"
	// "23P01": "exclusion_violation"
	case "23514", "23P01":
		errx.Code = CodeConditionNotMet
	}

	response := &ErrorResponse{
		StatusCode: http.StatusConflict,
		Err:        errx.Wrap(err),
	}

	return response
}

// PGDataError handles PG integrity errors
func PGDataError(err pq.Error) *ErrorResponse {
	errx := NewError(CodeConflict, "Data Error")

	switch err.Code {
	// "22003": "numeric_value_out_of_range",
	// "22008": "datetime_field_overflow",
	// "22015": "interval_field_overflow",
	// "22022": "indicator_overflow",
	// "22P01": "floating_point_exception",
	case "22003", "22008", "22015", "22022", "22P01":
		errx.Code = CodeOutOfrange
	// "22004": "null_value_not_allowed",
	// "22002": "null_value_no_indicator_parameter",
	case "22002", "22004":
		errx.Code = CodeConditionNotMet
	}

	response := &ErrorResponse{
		StatusCode: http.StatusUnprocessableEntity,
		Err:        errx.Wrap(err),
	}

	return response
}
