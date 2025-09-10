package v1

import (
	"log/slog"

	"github.com/alexliesenfeld/health"
	"github.com/linuxunsw/vote/backend/internal/api/v1/handlers"
	"github.com/linuxunsw/vote/backend/internal/api/v1/middleware"
	"github.com/linuxunsw/vote/backend/internal/config"
	"github.com/linuxunsw/vote/backend/internal/mailer"
	"github.com/linuxunsw/vote/backend/internal/store"

	"github.com/danielgtaylor/huma/v2"
)

// Register mounts all the API v1 routes using Huma groups and middleware.
func Register(api huma.API, logger *slog.Logger, store store.Store, mailer mailer.Mailer, checker health.Checker) {
	// Base group for all v1 routes
	v1 := huma.NewGroup(api, "/api/v1")

	huma.Get(api, "/health", handlers.GetHealth(checker))

	// == Authentication Routes ==
	// huma.Register(v1, huma.Operation{
	// 	OperationID: "auth-request-otp",
	// 	Method:      "POST",
	// 	Path:        "/auth/request-otp",
	// 	Summary:     "Request an OTP",
	// 	Tags:        []string{"Auth"},
	// }, handlers.RequestOTP(store, mailer))
	//
	// huma.Register(v1, huma.Operation{
	// 	OperationID: "auth-verify-otp",
	// 	Method:      "POST",
	// 	Path:        "/auth/verify-otp",
	// 	Summary:     "Verify an OTP to get a JWT",
	// 	Tags:        []string{"Auth"},
	// }, handlers.VerifyOTP(store, jwtSecret))

	// This group requires a valid JWT for all its routes.
	userRoutes := huma.NewGroup(v1)
	userRoutes.UseSimpleModifier(func(op *huma.Operation) {
		op.Security = []map[string][]string{
			{"cookieAuth": {}},
		}
	})
	authMiddleware := middleware.CookieAuthenticator(api, config.Load().JWT)
	userRoutes.UseMiddleware(authMiddleware)

	huma.Register(userRoutes, huma.Operation{
		OperationID: "submit-nomination",
		Method:      "PUT",
		Path:        "/nomination",
		Summary:     "Submit self-nomination",
		Description: "Creates a self-nomination, replacing an existing one",
		Tags:        []string{"Nominations"},
	}, handlers.SubmitNomination(logger, store))

	huma.Register(userRoutes, huma.Operation{
		OperationID: "get-nomination",
		Method:      "GET",
		Path:        "/nomination",
		Summary:     "Get self-nomination",
		Description: "Retrieves an existing self-nomination (if any)",
		Tags:        []string{"Nominations"},
	}, handlers.GetNomination(logger, store))

	// huma.Register(userRoutes, huma.Operation{
	// 	OperationID: "get-ballot",
	// 	Method:      "GET",
	// 	Path:        "/voting/ballot",
	// 	Summary:     "Get the voting ballot",
	// 	Tags:        []string{"Voting"},
	// }, handlers.GetBallot(store))
	//
	// huma.Register(userRoutes, huma.Operation{
	// 	OperationID: "submit-vote",
	// 	Method:      "POST",
	// 	Path:        "/voting/vote",
	// 	Summary:     "Submit or update a vote",
	// 	Tags:        []string{"Voting"},
	// }, handlers.SubmitVote(store))

	// == Admin Routes ==
	// This group requires a valid JWT AND admin privileges.
	// adminRoutes := huma.NewGroup(userRoutes, "/admin")
	// adminRoutes.UseMiddleware(middleware.RequireAdmin(api))
	//
	// huma.Register(adminRoutes, huma.Operation{
	// 	OperationID: "admin-update-election-state",
	// 	Method:      "PUT",
	// 	Path:        "/election/state",
	// 	Summary:     "Update the election state",
	// 	Tags:        []string{"Admin"},
	// }, handlers.UpdateElectionState(store))
	//
	// huma.Register(adminRoutes, huma.Operation{
	// 	OperationID: "admin-upload-members",
	// 	Method:      "POST",
	// 	Path:        "/members/upload",
	// 	Summary:     "Upload a CSV of members",
	// 	Tags:        []string{"Admin"},
	// }, handlers.UploadMembers(store))
	//
	// huma.Register(adminRoutes, huma.Operation{
	// 	OperationID: "admin-dashboard-stream",
	// 	Method:      "GET",
	// 	Path:        "/dashboard/stream",
	// 	Summary:     "Stream live updates for the admin dashboard",
	// 	Tags:        []string{"Admin"},
	// }, handlers.DashboardStream(store))
}
