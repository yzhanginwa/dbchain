package rest

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"

	"github.com/gorilla/mux"
)


// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/polls", storeName), createPollHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/polls/{%s}", storeName, "id"), showPollHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/user-polls/{%s}", storeName, "address"), showUserPollsHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/ballot/{%s}/{%s}", storeName, "id", "address"), showBallotHandler(cliCtx, storeName)).Methods("GET")
}
