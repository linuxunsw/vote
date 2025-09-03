package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"github.com/linuxunsw/vote/backend/internal/config"
)

// Define custom context keys to avoid collisions.
type contextKey string

const userContextKey = contextKey("user")

// UserClaims represents the data stored in the JWT.
type UserClaims struct {
	ZID     string `json:"sub"`
	IsAdmin bool   `json:"isAdmin"`
}

// Authenticator creates a Huma middleware to validate JWTs.
func Authenticator(api huma.API, cfg config.JWTConfig) func(ctx huma.Context, next func(ctx huma.Context)) {
	return func(ctx huma.Context, next func(ctx huma.Context)) {
		authHeader := ctx.Header("Authorization")
		if authHeader == "" {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Missing Authorization header")
			return
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || parts[0] != "Bearer" {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Invalid Authorization header format")
			return
		}
		tokenString := parts[1]

		// Convert the given secret into type jwt.Key (symmetric key)
		key, err := jwk.Import([]byte(cfg.Secret))
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusInternalServerError, "Unable to authenticate")
			return
		}

		// Parse and verify the token, validate it too
		token, err := jwt.ParseString(tokenString, jwt.WithKey(jwa.HS256(), key), jwt.WithValidate(true), jwt.WithIssuer(cfg.Issuer), jwt.WithPedantic(true))
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		// Extract custom claims into UserClaims.
		claims := &UserClaims{}
		token.Get("isAdmin", &claims.IsAdmin)
		claims.ZID, _ = token.Subject()

		// Put the claims into the request context and create a new huma.Context.
		newCtx := context.WithValue(ctx.Context(), userContextKey, claims)
		ctx = huma.WithContext(ctx, newCtx)

		next(ctx)
	}
}

// RequireAdmin is a middleware that checks for admin privileges from the context.
// It must be placed *after* the Authenticator middleware.
func RequireAdmin(api huma.API) func(ctx huma.Context, next func(ctx huma.Context)) {
	return func(ctx huma.Context, next func(ctx huma.Context)) {
		claims, ok := ctx.Context().Value(userContextKey).(*UserClaims)
		if !ok || !claims.IsAdmin {
			huma.WriteErr(api, ctx, http.StatusForbidden, "Admin access required")
			return
		}
		next(ctx)
	}
}

// GetUser retrieves user info from a plain context.Context.
func GetUser(ctx context.Context) (*UserClaims, bool) {
	info, ok := ctx.Value(userContextKey).(*UserClaims)
	return info, ok
}
