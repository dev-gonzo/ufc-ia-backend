package ufc

import "errors"

var (
	ErrInvalidAthleteSlug   = errors.New("invalid athlete slug")
	ErrAthleteNotFound      = errors.New("athlete not found")
	ErrUnexpectedStatusCode = errors.New("unexpected status code")
	ErrParseFailed          = errors.New("parse failed")
	ErrImageDecodeFailed    = errors.New("image decode failed")
	ErrImageEncodeFailed    = errors.New("image encode failed")
)
