package users

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

const (
	birthDateFormat string = "2006-01-02"
	passwordChars   string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	emailMaxLen     int    = 100
	minAge          int    = 18
	nameMaxLen      int    = 50
	passwordMaxLen  int    = 128
	passwordMinLen  int    = 8
)

var (
	emailRegex *regexp.Regexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	ErrBirthDateInvalidDate         error = errors.New("birth date cannot be in the future")
	ErrBirthDateInvalidFormat       error = errors.New("invalid birth date format. use YYYY-MM-DD")
	ErrBirthDateMinAge              error = fmt.Errorf("must be at least %d years old to register", minAge)
	ErrBirthDateRequired            error = errors.New("birth date is required")
	ErrEmailInvalidFormat           error = errors.New("invalid email address")
	ErrEmailRequired                error = errors.New("email is required")
	ErrEmailTooLong                 error = fmt.Errorf("email cannot be longer than %d characters", passwordMaxLen)
	ErrFirstNameRequired            error = errors.New("first name is required")
	ErrFirstNameTooLong             error = fmt.Errorf("first name cannot be longer than %d characters", nameMaxLen)
	ErrLastNameRequired             error = errors.New("last name is required")
	ErrLastNameTooLong              error = fmt.Errorf("last name cannot be longer than %d characters", nameMaxLen)
	ErrPasswordConfirmationMismatch error = errors.New("passwords do not match")
	ErrPasswordInvalidFormat        error = errors.New("password must contain at least 1 number, 1 uppercase and 1 lowercase letter")
	ErrPasswordMaxLen               error = fmt.Errorf("password cannot be longer than %d characters", passwordMinLen)
	ErrPasswordMinLen               error = fmt.Errorf("password must be at least %d characters long", passwordMinLen)
	ErrPasswordRequired             error = errors.New("password is required")
)

func validateFirstName(name string) error {
	if name == "" {
		return ErrFirstNameRequired
	}

	if len(name) > nameMaxLen {
		return ErrFirstNameTooLong
	}
	return nil
}

func validateLastName(name string) error {
	if name == "" {
		return ErrLastNameRequired
	}

	if len(name) > nameMaxLen {
		return ErrLastNameTooLong
	}
	return nil
}

func validateEmail(email string) error {
	if email == "" {
		return ErrEmailRequired
	}

	if len(email) > emailMaxLen {
		return ErrEmailTooLong
	}

	if !emailRegex.MatchString(email) {
		return ErrEmailInvalidFormat
	}

	return nil
}

func validateBirthDate(birthDate string) error {
	if birthDate == "" {
		return ErrBirthDateRequired
	}

	bday, err := time.Parse(birthDateFormat, birthDate)
	if err != nil {
		return ErrBirthDateInvalidFormat
	}

	if bday.After(time.Now()) {
		return ErrBirthDateInvalidDate
	}

	if int(time.Since(bday).Hours()/24/365) < minAge {
		return ErrBirthDateMinAge
	}

	return nil
}

func validatePassword(password, passwordConfirmation string) error {
	if password == "" {
		return ErrPasswordRequired
	}

	if len(password) < passwordMinLen {
		return ErrPasswordMinLen
	}

	if len(password) > passwordMaxLen {
		return ErrPasswordMaxLen
	}

	if !strings.ContainsAny(password, passwordChars) {
		return ErrPasswordInvalidFormat
	}

	if password != passwordConfirmation {
		return ErrPasswordConfirmationMismatch
	}
	return nil
}
