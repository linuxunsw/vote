package handlers

import (
	"context"
	"errors"
	"log/slog"

	"github.com/danielgtaylor/huma/v2"
	"github.com/linuxunsw/vote/backend/internal/api/v1/middleware/requestid"
	"github.com/linuxunsw/vote/backend/internal/api/v1/models"
	"github.com/linuxunsw/vote/backend/internal/store"
)

func ElectionMemberListSet(log *slog.Logger, st store.ElectionStore) func(ctx context.Context, input *models.ElectionMemberListSetInput) (*models.ElectionMemberListSetResponse, error) {
	return func(ctx context.Context, input *models.ElectionMemberListSetInput) (*models.ElectionMemberListSetResponse, error) {
		err := st.SetMembers(ctx, input.ElectionId, input.Body.Zids)
		if err != nil {
			log.Error("failed to set election members", "error", err, "request_id", requestid.Get(ctx))
			return nil, huma.Error500InternalServerError("internal error")
		}
		return &models.ElectionMemberListSetResponse{}, nil
	}
}

func CreateElection(log *slog.Logger, st store.ElectionStore) func(ctx context.Context, input *models.CreateElectionInput) (*models.CreateElectionResponse, error) {
	return func(ctx context.Context, input *models.CreateElectionInput) (*models.CreateElectionResponse, error) {
		electionId, err := st.CreateElection(ctx, input.Body.Name)
		if errors.Is(err, store.ErrElectionCreateAlreadyRunning) {
			log.Warn("attempted to create election while one is already running", "request_id", requestid.Get(ctx))
			return nil, huma.Error409Conflict("an election is already running")
		} else if err != nil {
			log.Error("failed to create election", "error", err, "request_id", requestid.Get(ctx))
			return nil, huma.Error500InternalServerError("internal error")
		}
		return &models.CreateElectionResponse{
			Body: models.CreateElectionResponseBody{
				ElectionId: electionId,
			},
		}, nil
	}
}
