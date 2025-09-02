package handlers

import (
	"context"

	"github.com/linuxunsw/vote/backend/internal/api/v1/models"
)

func CheckNominationStatus(store any) func(ctx context.Context, input *models.NominationStatusInput) (*models.NominationStatusResponse, error) {
	return func(ctx context.Context, input *models.NominationStatusInput) (*models.NominationStatusResponse, error) {
		panic("hi")
	}
}
