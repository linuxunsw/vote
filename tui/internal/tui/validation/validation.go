package validation

import (
	"errors"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
)

var (
	errEmail    = errors.New("please enter a valid email address")
	errOTP      = errors.New("please enter a valid verification code")
	errZID      = errors.New("please enter a valid zID")
	errRoles    = errors.New("please select a role")
	errURL      = errors.New("please enter a valid url (including the https://, etc.)")
	errNonEmpty = errors.New("this field cannot be empty")
)

// Validates a zID
// A zID must be in the following format: 'z5555555'
func ZID(zID string) error {
	var re = regexp.MustCompile("(?m)^z[0-9]{7}$")

	if !re.MatchString(zID) {
		return errZID
	}

	return nil
}

// Validates an OTP
// An OTP must be a 6 digit number
func OTP(otp string) error {
	var re = regexp.MustCompile("(?m)^[0-9]{6}$")

	if !re.MatchString(otp) {
		return errOTP
	}

	return nil
}

// Validates an email
func Email(email string) error {
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return errEmail
	}

	parts := strings.SplitN(addr.Address, "@", 2)
	if len(parts) != 2 || parts[1] == "" {
		return errEmail
	}

	// Check that the email domain is valid
	domain := strings.TrimSpace(parts[1])
	domain = strings.TrimSuffix(domain, ".")

	u, err := url.Parse("https://" + domain)
	if err != nil {
		return errEmail
	}
	if u.Host == "" {
		return errEmail
	}
	if !strings.Contains(u.Host, ".") {
		return errEmail
	}

	return nil
}

// Validates role selection
// Ensures the user must nominate for at least one role
func Role(roles []string) error {
	if len(roles) == 0 {
		return errRoles
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
		return errURL
	}

	return nil
}

// Validates any string field
// Makes the field mandatory - 'cannot be empty'
func NotEmpty(s string) error {
	if s == "" {
		return errNonEmpty
	}

	return nil
}
