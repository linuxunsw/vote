package handlers

import (
	"context"
	"log/slog"

	"github.com/danielgtaylor/huma/v2"
	"github.com/linuxunsw/vote/backend/internal/api/v1/models"
)

func SubmitNomination(logger *slog.Logger, store any) func(ctx context.Context, input *models.SubmitNominationRequest) (*models.SubmitNominationResponse, error) {
	return func(ctx context.Context, input *models.SubmitNominationRequest) (*models.SubmitNominationResponse, error) {
		logger.Info("Submit Nomination Info", "request", input)
		return nil, huma.Error500InternalServerError("stub")
	}
}

func GetNomination(logger *slog.Logger, store any) func(ctx context.Context, input *struct{}) (*models.GetNominationResponse, error) {
	return func(ctx context.Context, input *struct{}) (*models.GetNominationResponse, error) {
		logger.Info("Get Nomination Info", "request", input)
		return nil, huma.Error500InternalServerError("stub")
	}
}
