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

func SubmitBallot(log *slog.Logger, st store.BallotStore, el store.ElectionStore, nom store.NominationStore) func(ctx context.Context, input *models.SubmitBallotInput) (*struct{}, error) {
	return func(ctx context.Context, input *models.SubmitBallotInput) (*struct{}, error) {
		claims, valid := middleware.GetUser(ctx)
		if !valid {
			log.Warn("unauthenticated user tried to submit vote", "request_id", requestid.Get(ctx))
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

		// the keys of input.Body.Positions have been validated
		// the values have been validated for the correct zID format, but not if they are running 
		for position, candidateZid := range input.Body.Positions {
			nomination, err := nom.GetNomination(ctx, election.ElectionID, candidateZid)
			if err != nil {
				log.Error("failed to get nomination", "error", err, "zid", candidateZid, "election_id", election.ElectionID, "request_id", requestid.Get(ctx))
				return nil, huma.Error500InternalServerError("internal error")
			}
			if nomination == nil {
				log.Warn("candidate is not running", "zid", candidateZid, "position", position, "election_id", election.ElectionID, "request_id", requestid.Get(ctx))
				return nil, huma.Error400BadRequest("candidate is not running for position " + position)
			}
			if !nomination.IsRunningFor(position) {
				log.Warn("candidate is not running for position", "zid", candidateZid, "position", position, "election_id", election.ElectionID, "request_id", requestid.Get(ctx))
				return nil, huma.Error400BadRequest("candidate is not running for position " + position)
			}
		}

		err = st.SubmitOrReplaceBallot(ctx, election.ElectionID, claims.ZID, store.SubmitBallot{
			Positions: input.Body.Positions,
		})
		if err != nil {
			log.Error("failed to submit ballot", "error", err, "zid", claims.ZID, "election_id", election.ElectionID, "request_id", requestid.Get(ctx))
			return nil, huma.Error500InternalServerError("internal error")
		}
		return &struct{}{}, nil
	}
}

func GetBallot(log *slog.Logger, st store.BallotStore, el store.ElectionStore) func(ctx context.Context, input *struct {}) (*models.GetBallotResponse, error) {
	return func(ctx context.Context, input *struct {}) (*models.GetBallotResponse, error) {
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
		if ballot == nil {
			return nil, huma.Error404NotFound("no ballot found for user")
		}
		return &models.GetBallotResponse{
			Body: models.Ballot{
				Positions: ballot.Positions,
				CreatedAt: ballot.CreatedAt,
				UpdatedAt: ballot.UpdatedAt,
			},
		}, nil
	}
}

func DeleteBallot(log *slog.Logger, st store.BallotStore, el store.ElectionStore) func(ctx context.Context, input *struct {}) (*struct{}, error) {
	return func(ctx context.Context, input *struct {}) (*struct{}, error) {
		claims, valid := middleware.GetUser(ctx)
		if !valid {
			log.Warn("unauthenticated user tried to delete ballot", "request_id", requestid.Get(ctx))
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

		err = st.TryDeleteBallot(ctx, election.ElectionID, claims.ZID)
		if err != nil {
			log.Error("failed to delete ballot", "error", err, "zid", claims.ZID, "election_id", election.ElectionID, "request_id", requestid.Get(ctx))
			return nil, huma.Error500InternalServerError("internal error")
		}
		return &struct{}{}, nil
	}
}