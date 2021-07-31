package keeper_test

import (
	abci "github.com/dbchaincloud/tendermint/abci/types"

	"github.com/dbchaincloud/cosmos-sdk/simapp"
	sdk "github.com/dbchaincloud/cosmos-sdk/types"
	"github.com/dbchaincloud/cosmos-sdk/x/auth"
)

func createTestApp(isCheckTx bool) (*simapp.SimApp, sdk.Context) {
	app := simapp.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})

	app.AccountKeeper.SetParams(ctx, auth.DefaultParams())
	app.BankKeeper.SetSendEnabled(ctx, true)

	return app, ctx
}
