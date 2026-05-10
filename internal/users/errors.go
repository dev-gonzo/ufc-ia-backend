package users

import "errors"

var (
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidUsername = errors.New("invalid username")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidRole     = errors.New("invalid role")
	ErrEmailInUse      = errors.New("email already exists")
	ErrUsernameInUse   = errors.New("username already exists")
)
