package handlers

import (
	"context"
	"errors"
	"log/slog"
	"time"

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
			return nil, huma.Error400BadRequest("an election is already running")
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

func GetElectionState(log *slog.Logger, st store.ElectionStore) func(ctx context.Context, input *struct{}) (*models.GetElectionStateResponse, error) {
	return func(ctx context.Context, input *struct{}) (*models.GetElectionStateResponse, error) {
		election, err := st.CurrentElection(ctx)
		if err != nil {
			log.Error("failed to get current election", "error", err, "request_id", requestid.Get(ctx))
			return nil, huma.Error500InternalServerError("internal error")
		}
		if election == nil {
			return &models.GetElectionStateResponse{
				Body: models.GetElectionStateResponseBody{
					State: "NO_ELECTION",
				},
			}, nil
		}
		return &models.GetElectionStateResponse{
			Body: models.GetElectionStateResponseBody{
				ElectionId:     election.ElectionID,
				State:          string(election.State),
				StateCreatedAt: election.StateCreatedAt().Format(time.RFC3339),
			},
		}, nil
	}
}

func TransitionElectionState(log *slog.Logger, st store.ElectionStore) func(ctx context.Context, input *models.TransitionElectionStateInput) (*models.TransitionElectionStateResponse, error) {
	return func(ctx context.Context, input *models.TransitionElectionStateInput) (*models.TransitionElectionStateResponse, error) {
		err := st.CurrentElectionSetState(ctx, string(input.Body.State))
		if errors.Is(err, store.ErrElectionNotFound) {
			log.Warn("attempted to transition state with no election running", "request_id", requestid.Get(ctx))
			return nil, huma.Error400BadRequest("no election is currently running")
		} else if errors.Is(err, store.ErrElectionTransitionInvalidState) {
			// NOTE: this should never happen, as the OpenAPI schema should prevent invalid states
			log.Warn("attempted to transition election to invalid state", "state", input.Body.State, "request_id", requestid.Get(ctx))
			return nil, huma.Error422UnprocessableEntity("invalid state")
		} else if errors.Is(err, store.ErrElectionTransitionCannotRegress) {
			log.Warn("attempted to transition election state by regressing", "state", input.Body.State, "request_id", requestid.Get(ctx))
			return nil, huma.Error400BadRequest("cannot transition to a previous state")
		} else if errors.Is(err, store.ErrElectionTransitionCannotJump) {
			log.Warn("attempted to transition election state by jumping", "state", input.Body.State, "request_id", requestid.Get(ctx))
			return nil, huma.Error400BadRequest("cannot skip states when transitioning")
		} else if err != nil {
			log.Error("failed to transition election state", "error", err, "request_id", requestid.Get(ctx))
			return nil, huma.Error500InternalServerError("internal error")
		}

		return &models.TransitionElectionStateResponse{}, nil
	}
}
