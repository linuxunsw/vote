package v1

import (
	"log/slog"
	"net/http"

	"github.com/alexliesenfeld/health"
	"github.com/linuxunsw/vote/backend/internal/api/v1/handlers"
	"github.com/linuxunsw/vote/backend/internal/api/v1/middleware"
	"github.com/linuxunsw/vote/backend/internal/api/v1/models"
	"github.com/linuxunsw/vote/backend/internal/config"
	"github.com/linuxunsw/vote/backend/internal/mailer"
	"github.com/linuxunsw/vote/backend/internal/store"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/sse"
)

type HandlerDependencies struct {
	Logger *slog.Logger

	Cfg config.Config

	Mailer mailer.Mailer
	Checker health.Checker

	// Stores
	OtpStore store.OTPStore
	ElectionStore store.ElectionStore
	NominationStore store.NominationStore
}

// Register mounts all the API v1 routes using Huma groups and middleware.
func Register(api huma.API, deps HandlerDependencies) {
	// Base group for all v1 routes
	v1 := huma.NewGroup(api, "/api/v1")

	huma.Get(api, "/health", handlers.GetHealth(deps.Checker))

	// == Authentication Routes ==
	huma.Register(v1, huma.Operation{
		OperationID: "generate-otp",
		Method:      http.MethodPost,
		Path:        "/otp/generate",
		Summary:     "Generate an OTP code",
		Tags:        []string{"OTP"},
	}, handlers.GenerateOTP(deps.Logger, deps.OtpStore, deps.Mailer))

	huma.Register(v1, huma.Operation{
		OperationID: "submit-otp",
		Method:      http.MethodPost,
		Path:        "/otp/submit",
		Summary:     "Submit an OTP to enter a session",
		Tags:        []string{"OTP"},
	}, handlers.SubmitOTP(deps.Logger, deps.OtpStore, deps.ElectionStore, deps.Cfg.JWT, deps.Cfg.Admin))

	// This group requires a valid JWT for all its routes.
	userRoutes := huma.NewGroup(v1)
	userRoutes.UseSimpleModifier(func(op *huma.Operation) {
		op.Security = []map[string][]string{
			{"cookieAuth": {}},
		}
	})
	authMiddleware := middleware.CookieAuthenticator(api, deps.Cfg.JWT)
	userRoutes.UseMiddleware(authMiddleware)

	// state updates via SSE
	sse.Register(userRoutes, huma.Operation{
		OperationID: "state",
		Method:      http.MethodGet,
		Path:        "/state",
		Summary:     "Get State (SSE)",
		Description: "Gets current election state changes as server sent events.",
		Tags:        []string{"State"},
	}, map[string]any{
		"stateChange": &models.StateChangeEvent{},
	},
		handlers.GetState(deps.Logger))

	// nomination
	huma.Register(userRoutes, huma.Operation{
		OperationID: "submit-nomination",
		Method:      http.MethodPut,
		Path:        "/nomination",
		Summary:     "Submit self-nomination",
		Description: "Creates a self-nomination for the current election, replacing an existing one.",
		Tags:        []string{"Nominations"},
	}, handlers.SubmitNomination(deps.Logger, deps.NominationStore, deps.ElectionStore))

	huma.Register(userRoutes, huma.Operation{
		OperationID: "get-nomination",
		Method:      http.MethodGet,
		Path:        "/nomination",
		Summary:     "Get self-nomination",
		Description: "Retrieves an existing self-nomination (if any) for the current election.",
		Tags:        []string{"Nominations"},
	}, handlers.GetNomination(deps.Logger, deps.NominationStore, deps.ElectionStore))

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
	adminRoutes := huma.NewGroup(userRoutes)
	adminRoutes.UseMiddleware(middleware.RequireAdmin(api))
	
	huma.Register(adminRoutes, huma.Operation{
		OperationID: "create-election",
		Method:      http.MethodPost,
		Path:        "/elections",
		Summary:     "Create an election",
		Tags:        []string{"Elections"},
	}, handlers.CreateElection(deps.Logger, deps.ElectionStore))

	huma.Register(adminRoutes, huma.Operation{
		OperationID: "set-election-members",
		Method:      http.MethodPut,
		Path:        "/elections/{election_id}/members",
		Summary:     "Set the member list for an election",
		Tags:        []string{"Elections"},
	}, handlers.ElectionMemberListSet(deps.Logger, deps.ElectionStore))

	// huma.Register(adminRoutes, huma.Operation{
	// 	OperationID: "admin-update-election-state",
	// 	Method:      "PUT",
	// 	Path:        "/elections/state",
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
