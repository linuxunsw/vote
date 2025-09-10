package handlers

import (
	"context"

	"github.com/danielgtaylor/huma/v2/sse"
	"github.com/linuxunsw/vote/backend/internal/api/v1/models"
)

func GetState(ctx context.Context, input *struct{}, send sse.Sender) {
	// dummy code that sends random state change events
	// for x := 0; x < 10; x++ {
	// 	send.Data(models.StateChangeEvent{NewState: models.StateClosed.String()})
	// 	time.Sleep(1 * time.Second)
	// 	send.Data(models.StateChangeEvent{NewState: models.StateVotingOpen.String()})
	// 	time.Sleep(1 * time.Second)
	// }

	send.Data(models.StateChangeEvent{NewState: models.StateClosed.String()})
}
