package gov

import "errors"

var (
	ErrInvalidBytes  = errors.New("invalid bytes")
	ErrInvalidCoins  = errors.New("invalid coins")
	ErrInvalidString = errors.New("invalid string")
	ErrInvalidBool   = errors.New("invalid bool")
	ErrInvalidBigint = errors.New("invalid bigint")
)
