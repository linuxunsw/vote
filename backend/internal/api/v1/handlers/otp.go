package handlers

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"github.com/linuxunsw/vote/backend/internal/api/v1/middleware/requestid"
	"github.com/linuxunsw/vote/backend/internal/api/v1/models"
	"github.com/linuxunsw/vote/backend/internal/config"
	"github.com/linuxunsw/vote/backend/internal/mailer"
	"github.com/linuxunsw/vote/backend/internal/store"
)

// Generates a 6-digit OTP code
func NewCode() (string, error) {
	const maxOTP = 999999

	b := make([]byte, 4)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", fmt.Errorf("failed to read random bytes: %w", err)
	}

	num32 := int(b[0])<<24 | int(b[1])<<16 | int(b[2])<<8 | int(b[3])
	otpInt := num32 % (maxOTP + 1)

	otpString := strconv.Itoa(otpInt)
	for len(otpString) < 6 {
		otpString = "0" + otpString
	}

	return otpString, nil
}

func EmailFromZid(zid string) string {
	return zid + "@unsw.edu.au"
}

// Huma generate OTP handler
func GenerateOTP(log *slog.Logger, st store.OTPStore, mailer mailer.Mailer) func(ctx context.Context, input *models.GenerateOTPInput) (*models.GenerateOTPResponse, error) {
	return func(ctx context.Context, input *models.GenerateOTPInput) (*models.GenerateOTPResponse, error) {
		code, err := NewCode()
		if err != nil {
			log.Error("failed to generate OTP code", "error", err, "request_id", requestid.Get(ctx))
			return nil, huma.Error500InternalServerError("internal error")
		}

		err = st.CreateOrReplace(ctx, input.Body.Zid, code)
		if errors.Is(err, store.ErrOTPRateLimitExceeded) {
			log.Warn("rate limit exceeded for OTP generation", "zid", input.Body.Zid, "request_id", requestid.Get(ctx))
			return nil, huma.Error429TooManyRequests("rate limit exceeded")
		} else if err != nil {
			log.Error("failed to create OTP entry", "error", err, "request_id", requestid.Get(ctx))
			return nil, huma.Error500InternalServerError("internal error")
		}

		err = mailer.SendOTP(EmailFromZid(input.Body.Zid), code)
		if err != nil {
			log.Error("failed to send OTP email", "error", err, "request_id", requestid.Get(ctx))
			return nil, huma.Error500InternalServerError("internal error")
		}

		// success
		return &models.GenerateOTPResponse{}, nil
	}
}

func isAdminConfig(adminCfg config.AdminConfig, zid string) bool {
	for _, admin := range adminCfg.AdminZIds {
		if zid == admin {
			return true
		}
	}
	return false
}

// Huma submit OTP handler
func SubmitOTP(log *slog.Logger, st store.OTPStore, el store.ElectionStore, cfg config.JWTConfig, adminCfg config.AdminConfig) func(ctx context.Context, input *models.SubmitOTPInput) (*models.SubmitOTPResponse, error) {
	return func(ctx context.Context, input *models.SubmitOTPInput) (*models.SubmitOTPResponse, error) {
		valid, reason, err := st.ValidateAndConsume(ctx, input.Body.Zid, input.Body.Otp)
		if err != nil {
			log.Error("failed to validate OTP", "error", err, "request_id", requestid.Get(ctx))
			return nil, huma.Error500InternalServerError("internal error")
		}

		if !valid {
			clientStr := ""
			switch reason {
			case store.OTPValidateNotFoundOrExpired:
				clientStr = "invalid code"
			case store.OTPValidateAttemptsExceeded:
				clientStr = "attempts exceeded"
			case store.OTPValidateMismatch:
				clientStr = "invalid code"
			}

			log.Warn("invalid OTP submission", "zid", input.Body.Zid, "reason", reason.ToString(), "request_id", requestid.Get(ctx))
			return nil, huma.Error400BadRequest(clientStr)
		}

		isAdmin := isAdminConfig(adminCfg, input.Body.Zid)

		if !isAdmin {
			// admins need not be members, but even if we don't include this, you wouldn't be able to create elections
			currentElection, err := el.CurrentElection(ctx)
			if err != nil {
				log.Error("failed to get current election", "error", err, "request_id", requestid.Get(ctx))
				return nil, huma.Error500InternalServerError("internal error")
			}
			if currentElection == nil {
				log.Warn("no current election when submitting OTP", "zid", input.Body.Zid, "request_id", requestid.Get(ctx))
				return nil, huma.Error400BadRequest("no election is currently running")
			}

			entry, err := el.GetMember(ctx, currentElection.ElectionID, input.Body.Zid)
			if err != nil {
				log.Error("failed to get election member", "error", err, "request_id", requestid.Get(ctx))
				return nil, huma.Error500InternalServerError("internal error")
			}
			if entry == nil {
				log.Warn("zid not in election member list", "zid", input.Body.Zid, "request_id", requestid.Get(ctx))
				return nil, huma.Error403Forbidden("not authorized to vote")
			}
		} else {
			log.Info("admin login", "zid", input.Body.Zid, "request_id", requestid.Get(ctx))
		}

		tokenExpiry := time.Now().Add(cfg.Duration)

		// user is now authenticated and authorised as a society member
		// create JWT
		token, err := jwt.NewBuilder().
			Issuer(cfg.Issuer).
			Subject(input.Body.Zid).
			Claim("isAdmin", isAdmin).
			Expiration(tokenExpiry).
			Build()
		if err != nil {
			log.Error("failed to build JWT", "error", err, "request_id", requestid.Get(ctx))
			return nil, huma.Error500InternalServerError("internal error")
		}

		// sign JWT
		signed, err := jwt.Sign(token, jwt.WithKey(jwa.HS256(), []byte(cfg.Secret)))
		if err != nil {
			log.Error("failed to sign JWT", "error", err, "request_id", requestid.Get(ctx))
			return nil, huma.Error500InternalServerError("internal error")
		}

		// serialise into a string
		stringJWT := string(signed)

		return &models.SubmitOTPResponse{
			SetCookie: http.Cookie{
				Name:     cfg.CookieName,
				Value:    stringJWT,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
				Path:     "/",
				Expires:  tokenExpiry,
			},

			Body: models.SubmitOTPResponseBody{
				Zid:     input.Body.Zid,
				Expiry:  tokenExpiry,
				IsAdmin: isAdmin,
			},
		}, nil
	}
}
