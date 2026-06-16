package core

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTokenInvalid       = errors.New("token invalid")
	ErrInvitationInvalid  = errors.New("invitation invalid")
	ErrForbidden          = errors.New("forbidden")
	ErrProductNotFound    = errors.New("product not found")
	ErrAlreadyMember      = errors.New("user is already a member of this product")
	// Later add brute force prevention error
)
