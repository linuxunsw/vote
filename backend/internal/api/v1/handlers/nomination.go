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

func SubmitNomination(logger *slog.Logger, st store.NominationStore, el store.ElectionStore) func(ctx context.Context, input *models.SubmitNominationRequest) (*struct{}, error) {
	return func(ctx context.Context, input *models.SubmitNominationRequest) (*struct{}, error) {
		claims, valid := middleware.GetUser(ctx)
		if !valid {
			logger.Warn("unauthenticated user tried to submit nomination", "request_id", requestid.Get(ctx))
			return nil, huma.Error401Unauthorized("invalid user")
		}

		candidateZid := claims.ZID
		election, err := el.CurrentElection(ctx)
		if err != nil {
			logger.Error("failed to get current election", "error", err, "request_id", requestid.Get(ctx))
			return nil, huma.Error500InternalServerError("internal error")
		}
		if election == nil {
			logger.Warn("no current election when submitting OTP", "zid", candidateZid, "request_id", requestid.Get(ctx))
			return nil, huma.Error400BadRequest("no election is currently running")
		}

		err = st.SubmitOrReplaceNomination(ctx, election.ElectionID, candidateZid, input.Body)
		if err != nil {
			logger.Error("failed to submit or replace nomination", "error", err, "zid", candidateZid, "election_id", election.ElectionID, "request_id", requestid.Get(ctx))
			return nil, huma.Error500InternalServerError("internal error")
		}
		return &struct{}{}, nil
	}
}

func GetNomination(logger *slog.Logger, st store.NominationStore, el store.ElectionStore) func(ctx context.Context, input *struct{}) (*models.GetNominationResponse, error) {
	return func(ctx context.Context, input *struct{}) (*models.GetNominationResponse, error) {
		claims, valid := middleware.GetUser(ctx)
		if !valid {
			logger.Warn("unauthenticated user tried to submit nomination", "request_id", requestid.Get(ctx))
			return nil, huma.Error401Unauthorized("invalid user")
		}

		candidateZid := claims.ZID
		election, err := el.CurrentElection(ctx)
		if err != nil {
			logger.Error("failed to get current election", "error", err, "request_id", requestid.Get(ctx))
			return nil, huma.Error500InternalServerError("internal error")
		}
		if election == nil {
			logger.Warn("no current election when submitting OTP", "zid", candidateZid, "request_id", requestid.Get(ctx))
			return nil, huma.Error400BadRequest("no election is currently running")
		}

		nomination, err := st.GetNomination(ctx, election.ElectionID, candidateZid)
		if err != nil {
			logger.Error("failed to get nomination", "error", err, "zid", candidateZid, "election_id", election.ElectionID, "request_id", requestid.Get(ctx))
			return nil, huma.Error500InternalServerError("internal error")
		}
		if nomination == nil {
			return nil, huma.Error404NotFound("nomination not found")
		}
		return &models.GetNominationResponse{Body: *nomination}, nil
	}
}
