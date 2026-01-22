package errors

import "errors"

var (
	ErrNotFound         = errors.New("resource not found")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden")
	ErrInternalServer   = errors.New("internal server error")
	ErrBadRequest       = errors.New("bad request")
	ErrConflict         = errors.New("resource conflict")
	ErrValidationFailed = errors.New("validation failed")
)
