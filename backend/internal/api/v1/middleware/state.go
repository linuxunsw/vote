package middleware

import (
	"log/slog"
	"net/http"
	"slices"

	"github.com/danielgtaylor/huma/v2"
	"github.com/linuxunsw/vote/backend/internal/api/v1/middleware/requestid"
	"github.com/linuxunsw/vote/backend/internal/store"
)

type RequireElectionStateDeps struct {
	API           huma.API
	Logger        *slog.Logger
	ElectionStore store.ElectionStore
}

// limits access to the attached operation(s) based on the given allowed states.
// compares the current election's state with the allowlist, hence there must be a current
// election for a request to pass through this middleware
func RequireElectionState(deps RequireElectionStateDeps, allowedStates ...store.ElectionState) func(ctx huma.Context, next func(ctx huma.Context)) {
	return func(ctx huma.Context, next func(ctx huma.Context)) {
		el, err := deps.ElectionStore.CurrentElection(ctx.Context())
		if err != nil {
			deps.Logger.Error("Couldn't get current election in RequireElectionState", "error", err)
			_ = huma.WriteErr(deps.API, ctx, http.StatusInternalServerError, "Internal Error")
			return
		}

		if el == nil {
			deps.Logger.Warn("Attempt to request route with no current election", "request_id", requestid.Get(ctx.Context()), "path", ctx.Operation().Path)
			_ = huma.WriteErr(deps.API, ctx, http.StatusForbidden, "No Current Election")
			return
		}

		if !slices.Contains(allowedStates, el.State) {
			deps.Logger.Warn("Attempt to request route with invalid election state", "request_id", requestid.Get(ctx.Context()), "path", ctx.Operation().Path)
			_ = huma.WriteErr(deps.API, ctx, http.StatusForbidden, "Invalid Election State")
			return
		}

		next(ctx)
	}
}
