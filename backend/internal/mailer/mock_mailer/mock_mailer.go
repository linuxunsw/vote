package mock_mailer

import "github.com/linuxunsw/vote/backend/internal/mailer"

type MockMailer struct {
	otpEmails map[string]string
}

func NewMockMailer() mailer.Mailer {
	return &MockMailer{
		otpEmails: make(map[string]string),
	}
}

func (m *MockMailer) SendOTP(toEmail string, otpCode string) error {
	m.otpEmails[toEmail] = otpCode
	return nil
}

// Retrieve the most recent OTP code sent for an email address.
// Returns empty string if no OTP was sent
func (m *MockMailer) MockRetrieveOTP(toEmail string) string {
	return m.otpEmails[toEmail]
}
