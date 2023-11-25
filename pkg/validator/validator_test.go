package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatePhoneNumber(t *testing.T) {
	tests := []struct {
		name         string
		phoneNumber  string
		wantMessages []string
		wantValid    bool
	}{
		{
			name:        "invalid phone number",
			phoneNumber: "a",
			wantMessages: []string{
				"phone number must be numeric",
				"phone number must be 10-13 digits",
				"phone number must start with 62",
			},
			wantValid: false,
		},
		{
			name:        "number contains non-numeric character",
			phoneNumber: "628123456789a",
			wantMessages: []string{
				"phone number must be numeric",
			},
			wantValid: false,
		},
		{
			name:        "number is too short",
			phoneNumber: "62812",
			wantMessages: []string{
				"phone number must be 10-13 digits",
			},
			wantValid: false,
		},
		{
			name:        "number is too long",
			phoneNumber: "6281299231855678",
			wantMessages: []string{
				"phone number must be 10-13 digits",
			},
			wantValid: false,
		},
		{
			name:        "number does not start with 62",
			phoneNumber: "08123456789",
			wantMessages: []string{
				"phone number must start with 62",
			},
			wantValid: false,
		},
		{
			name:         "valid phone number",
			phoneNumber:  "628123456789",
			wantMessages: []string{},
			wantValid:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMessages, gotValid := ValidatePhoneNumber(tt.phoneNumber)
			assert.Equal(t, tt.wantValid, gotValid)
			assert.Equal(t, tt.wantMessages, gotMessages)
		})
	}
}

func TestValidateFullName(t *testing.T) {
	tests := []struct {
		name         string
		fullName     string
		wantMessages []string
		wantValid    bool
	}{
		{
			name:     "full name is empty",
			fullName: "",
			wantMessages: []string{
				"full name must be 3-60 characters",
			},
			wantValid: false,
		},
		{
			name:     "full name is too short",
			fullName: "ab",
			wantMessages: []string{
				"full name must be 3-60 characters",
			},
			wantValid: false,
		},
		{
			name:     "full name is too long",
			fullName: "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz",
			wantMessages: []string{
				"full name must be 3-60 characters",
			},
			wantValid: false,
		},
		{
			name:         "valid full name",
			fullName:     "John Doe",
			wantMessages: []string{},
			wantValid:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMessages, gotValid := ValidateFullName(tt.fullName)
			assert.Equal(t, tt.wantValid, gotValid)
			assert.Equal(t, tt.wantMessages, gotMessages)
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name         string
		password     string
		wantMessages []string
		wantValid    bool
	}{
		{
			name:     "password is empty",
			password: "",
			wantMessages: []string{
				"password must be 6-64 characters",
				"password must contain at least 1 uppercase letter, 1 number, and 1 special character",
			},
			wantValid: false,
		},
		{
			name:     "password is too short",
			password: "Ab9!",
			wantMessages: []string{
				"password must be 6-64 characters",
			},
			wantValid: false,
		},
		{
			name:     "password doesn't have uppercase letter",
			password: "ab9!ab",
			wantMessages: []string{
				"password must contain at least 1 uppercase letter, 1 number, and 1 special character",
			},
			wantValid: false,
		},
		{
			name:     "password doesn't have number",
			password: "Ab!abab",
			wantMessages: []string{
				"password must contain at least 1 uppercase letter, 1 number, and 1 special character",
			},
			wantValid: false,
		},
		{
			name:     "password doesn't have special character",
			password: "Ab9abab",
			wantMessages: []string{
				"password must contain at least 1 uppercase letter, 1 number, and 1 special character",
			},
			wantValid: false,
		},
		{
			name:         "valid password",
			password:     "Ab9!abab",
			wantMessages: []string{},
			wantValid:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMessages, gotValid := ValidatePassword(tt.password)
			assert.Equal(t, tt.wantValid, gotValid)
			assert.Equal(t, tt.wantMessages, gotMessages)
		})
	}
}
