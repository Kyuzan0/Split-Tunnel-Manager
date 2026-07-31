package domain

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrInvalidCIDR   = errors.New("invalid CIDR")
	ErrAlreadyExists = errors.New("already exists")
	ErrSecretMissing = errors.New("secret missing")
)
