package handlers

import (
	"context"
	"fmt"

	"github.com/linuxunsw/vote/backend/internal/api/v1/models"
	"github.com/linuxunsw/vote/backend/internal/config"
	"github.com/linuxunsw/vote/backend/internal/mailer"
	"github.com/linuxunsw/vote/backend/internal/store"
)

// Huma generate OTP handler
func GenerateOTP(st store.OTPStore, mailer mailer.Mailer) func(ctx context.Context, input *models.GenerateOTPInput) (*models.GenerateOTPResponse, error) {
	return func(ctx context.Context, input *models.GenerateOTPInput) (*models.GenerateOTPResponse, error) {
		fmt.Printf("%s", input.Body.Zid)
		panic("stub")
	}
}

// Huma submit OTP handler
func SubmitOTP(st store.OTPStore, cfg config.JWTConfig) func(ctx context.Context, input *models.SubmitOTPInput) (*models.SubmitOTPResponse, error) {
	return func(ctx context.Context, input *models.SubmitOTPInput) (*models.SubmitOTPResponse, error) {
		panic("stub")
	}
}
