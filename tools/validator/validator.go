package validator

// ValidatePhoneNumber validates phone number field based off certain criteria.
func ValidatePhoneNumber(phoneNumber string) (messages []string, valid bool) {
	messages = []string{}
	valid = true

	// phone number must be numeric
	for _, char := range phoneNumber {
		if char < '0' || char > '9' {
			messages = append(messages, "phone number must be numeric")
			valid = false
		}
	}

	// phone number must be 10-13 digits
	if len(phoneNumber) < 10 || len(phoneNumber) > 13 {
		messages = append(messages, "phone number must be 10-13 digits")
		valid = false
	}

	// phone number must start with 62
	if len(phoneNumber) < 2 || phoneNumber[0:2] != "62" {
		messages = append(messages, "phone number must start with 62")
		valid = false
	}

	return
}

// ValidateFullName validates full name field based off certain criteria.
func ValidateFullName(fullName string) (messages []string, valid bool) {
	messages = []string{}
	valid = true

	// full name must be 3-60 characters
	if len(fullName) < 3 || len(fullName) > 60 {
		messages = append(messages, "full name must be 3-60 characters")
		valid = false
	}

	return
}

// ValidatePassword validates password field based off certain criteria.
func ValidatePassword(password string) (messages []string, valid bool) {
	messages = []string{}
	valid = true

	// password must be 6-64 characters
	if len(password) < 6 || len(password) > 64 {
		messages = append(messages, "password must be 6-64 characters")
		valid = false
	}

	// password must contain at least 1 uppercase letter, 1 number, and 1 special character
	hasUppercase := false
	hasNumber := false
	hasSpecialChar := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUppercase = true
		case char >= '0' && char <= '9':
			hasNumber = true
		case char == '!' || char == '@' || char == '#' || char == '$' || char == '%' || char == '^' || char == '&':
			hasSpecialChar = true
		}
	}

	if !hasUppercase || !hasNumber || !hasSpecialChar {
		messages = append(messages, "password must contain at least 1 uppercase letter, 1 number, and 1 special character")
		valid = false
	}

	return
}
