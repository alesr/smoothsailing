package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"encore.dev/rlog"
	"encore.dev/storage/sqldb"
	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

var (
	db      *sqldb.Database = sqldb.Named("users")
	secrets struct{ jwtSecret string }
)

//encore:service
type Service struct {
	dbCtxTimeout     time.Duration
	jwtSigningMethod *jwt.SigningMethodHMAC
	jwtDuration      time.Duration
	jwtSigningKey    string
}

func initService() (*Service, error) {
	return &Service{
		dbCtxTimeout:     time.Second * 5,
		jwtSigningMethod: jwt.SigningMethodHS512,
		jwtDuration:      360 * time.Hour, // 15 days,
		jwtSigningKey:    secrets.jwtSecret,
	}, nil
}

// Register registers a new user.
//
//encore:api public method=POST path=/users
func (s *Service) Register(ctx context.Context, req *RegisterRequest) (*UserResponse, error) {
	req.normalize()

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("could not hash password: %s", err)
	}

	userID, err := gonanoid.New()
	if err != nil {
		return nil, ErrInternal
	}

	usr := User{
		ID:           userID,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		BirthDate:    req.BirthDate,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}

	dbCtx, cancel := context.WithTimeout(ctx, s.dbCtxTimeout)
	defer cancel()

	if err := storeUser(dbCtx, usr); err != nil {
		rlog.Error("could not store" + err.Error())
		return nil, ErrInternal
	}

	return &UserResponse{
		ID:        usr.ID,
		FirstName: usr.FirstName,
		LastName:  usr.LastName,
		Email:     usr.Email,
		BirthDate: usr.BirthDate,
		CreatedAt: usr.CreatedAt,
	}, nil
}

// Remove user by ID.
//
//encore:api public method=DELETE path=/users/:id
func (s *Service) Remove(ctx context.Context, id string) error {
	dbCtx, cancel := context.WithTimeout(ctx, s.dbCtxTimeout)
	defer cancel()

	if err := deleteUserByID(dbCtx, id); err != nil {
		rlog.Error("could not delete user by id" + err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return ErrInternal
	}
	return nil
}

// Fetch user by ID.
//
//encore:api public method=GET path=/users/:id
func (s *Service) FetchByID(ctx context.Context, id string) (*UserResponse, error) {
	dbCtx, cancel := context.WithTimeout(ctx, s.dbCtxTimeout)
	defer cancel()

	usr, err := getUserByID(dbCtx, id)
	if err != nil {
		rlog.Error("could not get user by id" + err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, ErrInternal
	}
	return &UserResponse{
		ID:        usr.ID,
		FirstName: usr.FirstName,
		LastName:  usr.LastName,
		Email:     usr.Email,
		BirthDate: usr.BirthDate,
		CreatedAt: usr.CreatedAt,
	}, nil
}

// List all registered users.
//
//encore:api private method=GET path=/users
func (s *Service) List(ctx context.Context) (*ListAllResponse, error) {
	dbCtx, cancel := context.WithTimeout(ctx, s.dbCtxTimeout)
	defer cancel()

	usrs, err := getAllUsers(dbCtx)
	if err != nil {
		rlog.Error("could not get all users" + err.Error())
		return nil, ErrInternal
	}

	usrsResp := make([]UserResponse, 0, len(usrs))

	for _, usr := range usrs {
		usrsResp = append(usrsResp, UserResponse{
			ID:        usr.ID,
			FirstName: usr.FirstName,
			LastName:  usr.LastName,
			Email:     usr.Email,
			BirthDate: usr.BirthDate,
			CreatedAt: usr.CreatedAt,
		})
	}

	return &ListAllResponse{
		Data: usrsResp,
	}, nil
}

// GetToken returns a JWT token for the user
//
//encore:api public method=POST path=/users/token
func (s *Service) GetToken(ctx context.Context, req *GetTokenRequest) (*GetTokenResponse, error) {
	usr, err := getByEmail(ctx, req.Email)
	if err != nil {
		rlog.Error("could not get user by email:" + err.Error())
		return nil, ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(usr.PasswordHash), []byte(req.Password)); err != nil {
		rlog.Error("missing token:" + err.Error())
		return nil, ErrUnauthorized
	}

	token, err := s.generateJWT(usr.ID)
	if err != nil {
		rlog.Error("missing token:" + err.Error())
		return nil, ErrUnauthorized
	}
	return &GetTokenResponse{
		Token: token,
	}, nil
}

// VerifyToken verifies a JWT token
func (s *Service) VerifyToken(ctx context.Context, token string) error {
	if token == "" {
		rlog.Error("missing token")
		return ErrTokenInvalid
	}

	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			rlog.Error("unexpected signing method: ", token.Header["alg"])
			return nil, ErrTokenInvalid
		}

		if method.Alg() != s.jwtSigningMethod.Alg() {
			rlog.Error("invalid token signing method", method.Alg(), s.jwtSigningMethod.Alg())
			return nil, errors.New("invalid token signing method")
		}
		return []byte(s.jwtSigningKey), nil
	})
	if err != nil {
		rlog.Error("could not parse token", err)
		return ErrTokenInvalid
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		rlog.Error("invalid token")
		return ErrTokenInvalid
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		rlog.Error("could not find user id in token")
		return ErrTokenInvalid
	}

	expiration, ok := claims["exp"].(float64)
	if !ok {
		rlog.Error("could not find expiration in token")
		return ErrTokenInvalid
	}

	if time.Unix(int64(expiration), 0).Before(time.Now()) {
		rlog.Error("token expired")
		return ErrTokenExpired
	}

	if _, err := getUserByID(ctx, userID); err != nil {
		rlog.Error("could not get user by id", err)
		return ErrTokenInvalid
	}
	return nil
}

func (s *Service) generateJWT(userID string) (string, error) {
	now := time.Now().UTC()

	token := jwt.NewWithClaims(s.jwtSigningMethod, jwtClaim{
		userID,
		jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(s.jwtDuration).Unix(),
		},
	})

	signedString, err := token.SignedString([]byte(s.jwtSigningKey))
	if err != nil {
		return "", fmt.Errorf("could not sign token: %s", err)
	}
	return signedString, nil
}
