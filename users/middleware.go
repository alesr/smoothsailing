package users

import (
	"encore.dev/beta/errs"
	"encore.dev/middleware"
)

//encore:middleware target=all
func ValidationMiddleware(req middleware.Request, next middleware.Next) middleware.Response {
	payload := req.Data().Payload
	if validator, ok := payload.(interface{ Validate() error }); ok {
		if err := validator.Validate(); err != nil {
			return middleware.Response{
				Err: errs.WrapCode(err, errs.InvalidArgument, "validation failed"),
			}
		}
	}
	return next(req)
}
