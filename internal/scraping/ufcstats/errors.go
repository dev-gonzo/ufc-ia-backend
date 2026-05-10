package ufcstats

import "errors"

var (
	ErrInvalidURL           = errors.New("invalid url")
	ErrRequestFailed        = errors.New("request failed")
	ErrUnexpectedStatusCode = errors.New("unexpected status code")
	ErrReadResponseBody     = errors.New("read response body failed")
	ErrParseFailed          = errors.New("parse failed")
	ErrMissingRequiredField = errors.New("missing required field")
)
