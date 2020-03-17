package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
)

// RegisterRoutes registers surprise-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.Use(mux.CORSMethodMiddleware(r))
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}
