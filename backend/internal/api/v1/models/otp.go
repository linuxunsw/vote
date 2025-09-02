package models

import "net/http"

type GenerateOTPInput struct {
	Body GenerateOTPInputBody
}

type GenerateOTPInputBody struct {
	Zid string `json:"zid" doc:"User zID" pattern:"^z[0-9]{7}$" example:"z0000000"`
}

type GenerateOTPResponse struct {
}

type SubmitOTPInput struct {
	Body SubmitOTPInputBody
}

type SubmitOTPInputBody struct {
	Zid string `json:"zid" doc:"User zID" pattern:"^z[0-9]{7}$" example:"z0000000"`
	Otp string `json:"otp" doc:"OTP Code" pattern:"^[0-9]{6}$" example:"123123"`
}

type SubmitOTPResponse struct {
	SetCookie http.Cookie `header:"Set-Cookie"`
}
