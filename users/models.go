package users

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// User is the domain model representing a user.
type User struct {
	ID           string
	FirstName    string
	LastName     string
	Email        string
	BirthDate    string
	PasswordHash string
	CreatedAt    time.Time
}

// UserResponse defines the user transport model.
type UserResponse struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	BirthDate string
	CreatedAt time.Time
}

// RegisterRequest represents the input payload necessary for registering a user.
type RegisterRequest struct {
	FirstName       string
	LastName        string
	Email           string
	BirthDate       string
	Password        string
	PasswordConfirm string
}

func (req *RegisterRequest) normalize() {
	caser := cases.Title(language.BrazilianPortuguese)
	req.FirstName = caser.String(strings.ToLower(req.FirstName))
	req.LastName = caser.String(strings.ToLower(req.LastName))
	req.Email = strings.ToLower(req.Email)
}

// Validate must be exported so it can be used by the validation middleware.
func (in *RegisterRequest) Validate() error {
	if err := validateFirstName(in.FirstName); err != nil {
		return err
	}

	if err := validateLastName(in.LastName); err != nil {
		return err
	}

	if err := validateEmail(in.Email); err != nil {
		return err
	}

	if err := validateBirthDate(in.BirthDate); err != nil {
		return err
	}

	if err := validatePassword(in.Password, in.PasswordConfirm); err != nil {
		return err
	}
	return nil
}

type ListAllResponse struct {
	Data []UserResponse
}

type GetTokenRequest struct {
	Email    string
	Password string
}

func (req *GetTokenRequest) Validate() error {
	if req.Email == "" {
		return ErrEmailRequired
	}

	if req.Password == "" {
		return ErrPasswordRequired
	}
	return nil
}

type GetTokenResponse struct {
	Token string
}

type jwtClaim struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}
