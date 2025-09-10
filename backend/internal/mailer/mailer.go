package mailer

type Mailer interface{
	SendOTP(toEmail string, otpCode string) error
}
