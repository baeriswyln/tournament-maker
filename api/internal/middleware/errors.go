// Brief   :
// Author  : Nicolas Baeriswyl <nicolas@baeriswyl.dev>
// Created : 14.03.2026

package middleware

import "errors"

// Sentinel errors returned by the auth package. Callers may compare against
// these with errors.Is to produce tailored HTTP responses.
var (
	// ErrMissingToken is returned when no Authorization header is present.
	ErrMissingToken = errors.New("auth: missing token")

	// ErrMalformedToken is returned when the Authorization header is not in
	// "Bearer <token>" format.
	ErrMalformedToken = errors.New("auth: malformed token")

	// ErrInvalidToken is returned when the token signature or structure is invalid.
	ErrInvalidToken = errors.New("auth: invalid token")

	// ErrTokenExpired is returned when a valid token has passed its expiry time.
	ErrTokenExpired = errors.New("auth: token expired")

	// ErrUnexpectedSigningMethod is returned when the token header names a
	// signing algorithm other than HMAC (e.g. RS256 confusion attacks).
	ErrUnexpectedSigningMethod = errors.New("auth: unexpected signing method")

	// ErrUserInvalid is returned when the UserValidatorFunc rejects the user,
	// e.g. the account was deleted or suspended since the token was issued.
	ErrUserInvalid = errors.New("auth: user is not valid")
)
