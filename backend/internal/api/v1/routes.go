package v1

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

// TODO: implement
func addRoutes(api huma.API) {
	huma.Get(api, "/health", func(ctx context.Context, i *struct{}) (*struct{}, error) {
		return nil, nil
	})
}
