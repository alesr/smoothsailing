package users

import (
	"context"
	"time"

	"encore.dev/beta/errs"
	"encore.dev/rlog"
	"encore.dev/storage/sqldb"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

var (
	db           *sqldb.Database = sqldb.Named("users")
	dbCtxTimeout time.Duration   = time.Second * 5
)

//encore:service
type Service struct{}

func initService() (*Service, error) {
	return &Service{}, nil
}

// Register registers a new user
//
//encore:api public method=POST path=/users
func (s *Service) Register(ctx context.Context, params *RegisterParams) (*RegisterResponse, error) {
	rlog.Debug("user registration request received", params)

	userID, err := gonanoid.New()
	if err != nil {
		return nil, errs.WrapCode(err, errs.Internal, "could not generate user id")
	}

	usr := User{
		id:        userID,
		firstName: params.FirstName,
		lastName:  params.LastName,
		email:     params.Email,
		birthDate: params.BirthDate,
		createdAt: time.Now(),
	}

	dbCtx, cancel := context.WithTimeout(ctx, dbCtxTimeout)
	defer cancel()

	if err := Insert(dbCtx, usr); err != nil {
		return nil, errs.WrapCode(err, errs.Internal, "could not store user")
	}

	return &RegisterResponse{
		ID: usr.id,
	}, nil
}
