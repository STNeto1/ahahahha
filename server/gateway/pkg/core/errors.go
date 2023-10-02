package core

import "errors"

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserDoesNotExists  = errors.New("user does not exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInternalError      = errors.New("internal error")
)
