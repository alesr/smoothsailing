package users

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterParamsNormalize(t *testing.T) {
	testCases := []struct {
		name     string
		given    *RegisterRequest
		expected *RegisterRequest
	}{
		{
			name: "nothing changes",
			given: &RegisterRequest{
				FirstName:       "Joe",
				LastName:        "Doe",
				Email:           "joedoe@bar.foo",
				BirthDate:       "1999-01-01",
				Password:        "joedoe123",
				PasswordConfirm: "joedoe123",
			},
			expected: &RegisterRequest{
				FirstName:       "Joe",
				LastName:        "Doe",
				Email:           "joedoe@bar.foo",
				BirthDate:       "1999-01-01",
				Password:        "joedoe123",
				PasswordConfirm: "joedoe123",
			},
		},
		{
			name: "names get capitalized and email lowered",
			given: &RegisterRequest{
				FirstName:       "JOE",
				LastName:        "DOE",
				Email:           "JOEDOE@BAR.FOO",
				BirthDate:       "1999-01-01",
				Password:        "joedoe123",
				PasswordConfirm: "joedoe123",
			},
			expected: &RegisterRequest{
				FirstName:       "Joe",
				LastName:        "Doe",
				Email:           "joedoe@bar.foo",
				BirthDate:       "1999-01-01",
				Password:        "joedoe123",
				PasswordConfirm: "joedoe123",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.given.normalize()
			assert.Equal(t, tc.expected, tc.given)
		})
	}
}
