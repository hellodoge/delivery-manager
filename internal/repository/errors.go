package repository

import "errors"

var (
	ErrUserExists           = errors.New("user already exists")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
)
