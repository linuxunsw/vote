package messages

// Sends the user's form input (their zID) to the root model
type SendAuthMsg struct {
	ZID string
}

// Triggers the root model to check if otp is valid
type CheckOTPMsg struct {
	OTP string
}

// Empty message for now unless data needs to be sent to UI
// When UI recieves this message we know we can move on in the form
type AuthenticatedMsg struct{}
