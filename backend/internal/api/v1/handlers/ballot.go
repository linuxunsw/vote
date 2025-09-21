package handlers

import (
	"context"
	"log/slog"

	"github.com/danielgtaylor/huma/v2"
	"github.com/linuxunsw/vote/backend/internal/api/v1/middleware"
	"github.com/linuxunsw/vote/backend/internal/api/v1/middleware/requestid"
	"github.com/linuxunsw/vote/backend/internal/api/v1/models"
	"github.com/linuxunsw/vote/backend/internal/store"
)

// Note that "ballots" as refered to in the backend/database actually mean "votes" in the frontend/handlers.
// This GetBallot is actually just returning the user's submitted vote.
func GetBallot(log *slog.Logger, st store.BallotStore, el store.ElectionStore, nom store.NominationStore) func(ctx context.Context, input *struct{}) (*models.GetBallotResponse, error) {
	return func(ctx context.Context, input *struct{}) (*models.GetBallotResponse, error) {
		claims, valid := middleware.GetUser(ctx)
		if !valid {
			log.Warn("unauthenticated user tried to get ballot", "request_id", requestid.Get(ctx))
			return nil, huma.Error401Unauthorized("invalid user")
		}

		election, err := el.CurrentElection(ctx)
		if err != nil {
			log.Error("failed to get current election", "error", err, "request_id", requestid.Get(ctx))
			return nil, huma.Error500InternalServerError("internal error")
		}
		if election == nil {
			log.Warn("no current election", "request_id", requestid.Get(ctx))
			return nil, huma.Error400BadRequest("no election is currently running")
		}

		ballot, err := st.GetBallot(ctx, election.ElectionID, claims.ZID)
		if err != nil {
			log.Error("failed to get ballot", "error", err, "zid", claims.ZID, "election_id", election.ElectionID, "request_id", requestid.Get(ctx))
			return nil, huma.Error500InternalServerError("internal error")
		}
		// this is all we use GetBallot for
		hasVoted := ballot != nil

		// executive role -> candidates running for that role
		candidates := make(map[string][]models.PublicNomination)

		nominations, err := nom.GetElectionNominations(ctx, election.ElectionID)
		if err != nil {
			log.Error("failed to get nominations", "error", err, "election_id", election.ElectionID, "request_id", requestid.Get(ctx))
			return nil, huma.Error500InternalServerError("internal error")
		}
		for _, nomination := range nominations {
			for _, role := range nomination.ExecutiveRoles {
				candidates[role] = append(candidates[role], models.FromStoreNomination(nomination))
			}
		}
		response := models.PublicBallot{
			ElectionID: election.ElectionID,
			HasVoted:   hasVoted,
			Candidates: candidates,
		}
		return &models.GetBallotResponse{Body: response}, nil
	}
}
