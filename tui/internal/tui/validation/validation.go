package validation

import (
	"errors"
	"net/mail"
	"net/url"
	"regexp"
)

// Validates a zID
// A zID must be in the following format: 'z5555555'
func ZID(zID string) error {
	var re = regexp.MustCompile("(?m)^z[0-9]{7}$")

	if !re.MatchString(zID) {
		return errors.New("please enter a valid zID")
	}

	return nil
}

// Validates an OTP
// An OTP must be a 6 digit number
func OTP(otp string) error {
	var re = regexp.MustCompile("(?m)^[0-9]{6}$")

	if !re.MatchString(otp) {
		return errors.New("please enter a valid verification code")
	}

	return nil
}

// FIX: improve email validation to require a correct domain
// e.g. hi@hi should not be a valid email
// Validates an email
func Email(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("please enter a valid email address")
	}

	return nil
}

// Validates role selection
// Ensures the user must nominate for at least one role
func Role(roles []string) error {
	if len(roles) == 0 {
		return errors.New("please select a role")
	}

	return nil
}

// Validates URL
// Ensures that if a URL is provided, it will be a valid url
func URL(formURL string) error {
	if formURL == "" {
		return nil
	}

	parsed, err := url.Parse(formURL)
	if err != nil || !parsed.IsAbs() {
		return errors.New("please enter a valid URL")
	}

	return nil
}

// Validates any string field
// Makes the field mandatory - 'cannot be empty'
func NotEmpty(s string) error {
	if s == "" {
		return errors.New("this field cannot be empty")
	}

	return nil
}
