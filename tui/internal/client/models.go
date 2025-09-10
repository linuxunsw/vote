package client

type generateOTPRequest struct {
	ZID string `json:"zid"`
}

type submitOTPRequest struct {
	OTP string `json:"otp"`
	ZID string `json:"zid"`
}

type submitNominationResponse struct {
	ID string `json:"id"`
}

type responseError struct {
	Errors []problemError `json:"errors"`
}

type problemError struct {
	Location string `json:"location"`
	Message  string `json:"message"`
	Value    any    `json:"value"`
}
