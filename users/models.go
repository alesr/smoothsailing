package users

import "time"

// User is the domain model representing a user
type User struct {
	id        string
	firstName string
	lastName  string
	email     string
	birthDate string
	createdAt time.Time
}

// RegisterParams represents the necessary payload for registering a user.
type RegisterParams struct {
	FirstName       string
	LastName        string
	Email           string
	BirthDate       string
	Password        string
	PasswordConfirm string
}

func (r *RegisterParams) Validate() error {
	if err := validateFirstName(r.FirstName); err != nil {
		return err
	}

	if err := validateLastName(r.LastName); err != nil {
		return err
	}

	if err := validateEmail(r.Email); err != nil {
		return err
	}

	if err := validateBirthDate(r.BirthDate); err != nil {
		return err
	}

	if err := validatePassword(r.Password, r.PasswordConfirm); err != nil {
		return err
	}
	return nil
}

// RegisterResponse represents the returned payload when a user is created.
type RegisterResponse struct {
	ID string // User ID
}
