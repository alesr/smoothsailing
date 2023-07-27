package users

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	given := RegisterRequest{
		FirstName:       "Joe",
		LastName:        "Doe",
		Email:           "joedoe@foo.bar",
		BirthDate:       "2000-01-01",
		Password:        "abcde1234",
		PasswordConfirm: "abcde1234",
	}

	expected := UserResponse{
		FirstName: "Joe",
		LastName:  "Doe",
		Email:     "joedoe@foo.bar",
		BirthDate: "2000-01-01",
	}

	svc := Service{
		dbCtxTimeout:     time.Second * 5,
		jwtSigningMethod: jwt.SigningMethodHS512,
		jwtDuration:      time.Minute,
		jwtSigningKey:    "secret",
	}

	observed, err := svc.Register(context.TODO(), &given)
	require.NoError(t, err)

	// Assert API response

	// Unpredictable
	assert.NotEmpty(t, observed.ID)
	assert.NotEmpty(t, observed.CreatedAt)

	assert.Equal(t, expected.FirstName, observed.FirstName)
	assert.Equal(t, expected.LastName, observed.LastName)
	assert.Equal(t, expected.Email, observed.Email)
	assert.Equal(t, expected.BirthDate, observed.BirthDate)

	// Assert DB entry

	observedStore, err := getByEmail(context.TODO(), observed.Email)
	require.NoError(t, err)

	assert.Equal(t, observed.ID, observedStore.ID)
	assert.Equal(t, observed.CreatedAt.UTC(), observedStore.CreatedAt)

	assert.Equal(t, expected.FirstName, observedStore.FirstName)
	assert.Equal(t, expected.LastName, observedStore.LastName)
	assert.Equal(t, expected.Email, observedStore.Email)
	assert.Equal(t, expected.BirthDate, observedStore.BirthDate)
}
