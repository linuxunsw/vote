package handlers

import (
	"context"
	"log/slog"

	"github.com/danielgtaylor/huma/v2"
	"github.com/linuxunsw/vote/backend/internal/api/v1/models"
	"github.com/linuxunsw/vote/backend/internal/store"
)

func SubmitNomination(logger *slog.Logger, store store.Store) func(ctx context.Context, input *models.NominationRequest) (*models.NominationResponse, error) {
	return func(ctx context.Context, input *models.NominationRequest) (*models.NominationResponse, error) {
		logger.Info("Submit Nomination Info", "request", input)
		return nil, huma.Error500InternalServerError("stub")
	}
}
