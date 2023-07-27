package users

import (
	"encore.dev/beta/errs"
)

var (
	// Enumerate API errors

	ErrInternal *errs.Error = &errs.Error{
		Code:    errs.Internal,
		Message: "an internal error has happened",
	}

	ErrUserNotFound *errs.Error = &errs.Error{
		Code:    errs.NotFound,
		Message: "could not find user",
	}

	ErrInvalidPassword *errs.Error = &errs.Error{
		Code:    errs.InvalidArgument,
		Message: "invalid password",
	}

	ErrTokenInvalid *errs.Error = &errs.Error{
		Code:    errs.Unauthenticated,
		Message: "invalid authentication",
	}

	ErrTokenExpired *errs.Error = &errs.Error{
		Code:    errs.Unauthenticated,
		Message: "token expired",
	}

	ErrUnauthorized *errs.Error = &errs.Error{
		Code:    errs.InvalidArgument,
		Message: "unauthorized",
	}
)
