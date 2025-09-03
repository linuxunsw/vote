package requestid

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type ctxKey struct{}

func ctxKeyVal() ctxKey { return ctxKey{} }

// HumaMiddleware returns a router-agnostic huma middleware that ensures
// each request has an ID. Provide a header name (defaults to "X-Request-ID" for empty string)
func HumaMiddleware(header string) func(ctx huma.Context, next func(huma.Context)) {
	if header == "" {
		header = "X-Request-ID"
	}

	return func(ctx huma.Context, next func(huma.Context)) {
		id := ctx.Header(header)
		if id == "" {
			id = uuid.NewString()
		}

		// set response header and put id into request context
		ctx.SetHeader(header, id)
		ctx = huma.WithValue(ctx, ctxKeyVal(), id)

		next(ctx)
	}
}

// Get extracts the request id from a context.Context (from HumaMiddleware).
func Get(ctx context.Context) string {
	if v, _ := ctx.Value(ctxKeyVal()).(string); v != "" {
		return v
	}
	return ""
}
