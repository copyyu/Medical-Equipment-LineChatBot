package errors

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAdminNotFound      = errors.New("admin not found")
	ErrAdminInactive      = errors.New("admin is inactive")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrUsernameExists     = errors.New("username already exists")
	ErrEmailExists        = errors.New("email already exists")
	ErrSessionExpired     = errors.New("session expired")
	ErrWeakPassword       = errors.New("password is too weak")
)
