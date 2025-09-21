package v1

import (
	"log/slog"
	"net/http"

	"github.com/alexliesenfeld/health"
	"github.com/linuxunsw/vote/backend/internal/api/v1/handlers"
	"github.com/linuxunsw/vote/backend/internal/api/v1/middleware"
	"github.com/linuxunsw/vote/backend/internal/config"
	"github.com/linuxunsw/vote/backend/internal/mailer"
	"github.com/linuxunsw/vote/backend/internal/store"

	"github.com/danielgtaylor/huma/v2"
)

type HandlerDependencies struct {
	Logger *slog.Logger

	Cfg config.Config

	Mailer  mailer.Mailer
	Checker health.Checker

	// Stores
	OtpStore        store.OTPStore
	ElectionStore   store.ElectionStore
	NominationStore store.NominationStore
	BallotStore     store.BallotStore
}

// Register mounts all the API v1 routes using Huma groups and middleware.
func Register(api huma.API, deps HandlerDependencies) {
	// Base group for all v1 routes
	v1 := huma.NewGroup(api, "/api/v1")

	huma.Get(api, "/health", handlers.GetHealth(deps.Checker))

	// deps for election state middleware
	esDeps := middleware.RequireElectionStateDeps{
		API:           api,
		Logger:        deps.Logger,
		ElectionStore: deps.ElectionStore,
	}

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

	// == Authenticated Routes ==
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
	/* sse.Register(userRoutes, huma.Operation{
		OperationID: "state",
		Method:      http.MethodGet,
		Path:        "/state",
		Summary:     "Get State (SSE)",
		Description: "Gets current election state changes as server sent events.",
		Tags:        []string{"State"},
	}, map[string]any{
		"stateChange": &models.StateChangeEvent{},
	},
		handlers.GetState(deps.Logger)) */

	// election state
	huma.Register(userRoutes, huma.Operation{
		OperationID: "get-election-state",
		Method:      "GET",
		Path:        "/state",
		Summary:     "Get the current election state",
		Tags:        []string{"State"},
	}, handlers.GetElectionState(deps.Logger, deps.ElectionStore))

	// == Nomination Routes ==
	// Group for nomination operations that REQUIRE nominations to be open
	nominationRoutes := huma.NewGroup(userRoutes)
	nominationRoutes.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"Nominations"}
	})
	nominationRoutes.UseMiddleware(middleware.RequireElectionState(esDeps, store.StateNominationsOpen))

	huma.Register(nominationRoutes, huma.Operation{
		OperationID: "submit-nomination",
		Method:      http.MethodPut,
		Path:        "/nomination",
		Summary:     "Submit self-nomination",
		Description: "Creates a self-nomination for the current election, replacing an existing one.",
	}, handlers.SubmitNomination(deps.Logger, deps.NominationStore, deps.ElectionStore))

	huma.Register(nominationRoutes, huma.Operation{
		OperationID: "delete-nomination",
		Method:      http.MethodDelete,
		Path:        "/nomination",
		Summary:     "Delete self-nomination",
		Description: "Deletes an existing self-nomination for the current election. If an election is running, this route will always return as it succeeded even if a nomination did not exist.",
	}, handlers.DeleteNomination(deps.Logger, deps.NominationStore, deps.ElectionStore))

	// Nomination read operations (no state restriction needed)
	huma.Register(userRoutes, huma.Operation{
		OperationID: "get-nomination",
		Method:      http.MethodGet,
		Path:        "/nomination",
		Summary:     "Get self-nomination",
		Description: "Retrieves an existing self-nomination (if any) for the current election.",
		Tags:        []string{"Nominations"},
	}, handlers.GetNomination(deps.Logger, deps.NominationStore, deps.ElectionStore))

	huma.Register(userRoutes, huma.Operation{
		OperationID: "get-public-nomination",
		Method:      "GET",
		Path:        "/nomination/{nomination_id}",
		Summary:     "Get a public nomination by ID",
		Tags:        []string{"Nominations"},
	}, handlers.GetPublicNomination(deps.Logger, deps.NominationStore))

	// == Voting Routes ==
	// Group for voting operations that REQUIRE voting to be open
	votingRoutes := huma.NewGroup(userRoutes)
	votingRoutes.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"Voting"}
	})
	votingRoutes.UseMiddleware(middleware.RequireElectionState(esDeps, store.StateVotingOpen))

	huma.Register(votingRoutes, huma.Operation{
		OperationID: "submit-vote",
		Method:      "PUT",
		Path:        "/vote",
		Summary:     "Submit or update your current vote",
	}, handlers.SubmitVote(deps.Logger, deps.BallotStore, deps.ElectionStore, deps.NominationStore))

	huma.Register(votingRoutes, huma.Operation{
		OperationID: "delete-vote",
		Method:      "DELETE",
		Path:        "/vote",
		Summary:     "Delete your current vote",
	}, handlers.DeleteVote(deps.Logger, deps.BallotStore, deps.ElectionStore))

	// Voting read operations (no state restriction needed)
	huma.Register(userRoutes, huma.Operation{
		OperationID: "get-vote",
		Method:      "GET",
		Path:        "/vote",
		Summary:     "Get your current vote",
		Tags:        []string{"Voting"},
	}, handlers.GetVote(deps.Logger, deps.BallotStore, deps.ElectionStore))

	huma.Register(userRoutes, huma.Operation{
		OperationID: "get-ballot",
		Method:      "GET",
		Path:        "/ballot",
		Summary:     "Get the current ballot",
		Tags:        []string{"Voting"},
	}, handlers.GetBallot(deps.Logger, deps.BallotStore, deps.ElectionStore, deps.NominationStore))

	// == Admin Routes ==
	// This group requires a valid JWT AND admin privileges.
	adminRoutes := huma.NewGroup(userRoutes)
	adminRoutes.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"Admin"}
	})
	adminRoutes.UseMiddleware(middleware.RequireAdmin(api))

	huma.Register(adminRoutes, huma.Operation{
		OperationID: "create-election",
		Method:      http.MethodPost,
		Path:        "/elections",
		Summary:     "Create an election",
	}, handlers.CreateElection(deps.Logger, deps.ElectionStore))

	huma.Register(adminRoutes, huma.Operation{
		OperationID: "set-election-members",
		Method:      http.MethodPut,
		Path:        "/elections/{election_id}/members",
		Summary:     "Set the member list for an election",
	}, handlers.ElectionMemberListSet(deps.Logger, deps.ElectionStore))

	huma.Register(adminRoutes, huma.Operation{
		OperationID: "admin-transition-election-state",
		Method:      "PUT",
		Path:        "/state",
		Summary:     "Transition the election state",
	}, handlers.TransitionElectionState(deps.Logger, deps.ElectionStore))

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
