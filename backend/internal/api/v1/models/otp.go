package models

import (
	"net/http"
	"time"
)

type GenerateOTPInput struct {
	Body struct {
		Zid string `json:"zid" doc:"User zID" pattern:"^z[0-9]{7}$" example:"z0000000"`
	}
}

type GenerateOTPResponse struct {
}

type SubmitOTPInput struct {
	Body struct {
		Zid string `json:"zid" doc:"User zID" pattern:"^z[0-9]{7}$" example:"z0000000"`
		Otp string `json:"otp" doc:"OTP Code" pattern:"^[0-9]{6}$" example:"123123"`
	}
}

type SubmitOTPResponse struct {
	SetCookie http.Cookie `header:"Set-Cookie"`

	Body SubmitOTPResponseBody
}

type SubmitOTPResponseBody struct {
	Zid     string    `json:"zid" doc:"User zID" pattern:"^z[0-9]{7}$" example:"z0000000"`
	Expiry  time.Time `json:"expiry" format:"date-time" example:"2024-01-15T10:30:00Z" doc:"Timestamp when your session expires."`
	IsAdmin bool      `json:"is_admin" doc:"Admin status"`
}
