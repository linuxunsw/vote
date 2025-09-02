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

func SubmitNomination(store any) func(ctx context.Context, input *models.NominationSubmissionInput) (*models.NominationSubmissionResponse, error) {
	return func(ctx context.Context, input *models.NominationSubmissionInput) (*models.NominationSubmissionResponse, error) {
		panic("hi")
	}
}

func DeleteNomination(store any) func(ctx context.Context, input *models.NominationDeleteInput) (*models.NominationDeleteResponse, error) {
	return func(ctx context.Context, input *models.NominationDeleteInput) (*models.NominationDeleteResponse, error) {
		panic("hi")
	}
}
