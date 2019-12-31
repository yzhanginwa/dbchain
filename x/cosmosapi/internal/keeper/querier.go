package keeper

import (
    "fmt"
    "github.com/cosmos/cosmos-sdk/codec"

    sdk "github.com/cosmos/cosmos-sdk/types"
    abci "github.com/tendermint/tendermint/abci/types"
    "github.com/yzhanginwa/cosmos-api/x/cosmosapi/internal/types"
)

// query endpoints supported by the cosmosapi service Querier
const (
    QueryTable   = "table"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
    return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
        switch path[0] {
        case QueryTable:
            return queryTable(ctx, path[1:], req, keeper)
        default:
            return nil, sdk.ErrUnknownRequest("unknown cosmosapi query endpoint")
        }
    }
}

func queryTable(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
    table, err := keeper.GetTable(ctx, path[0])

    if err != nil {
        return []byte{}, sdk.ErrUnknownRequest(fmt.Sprintf("Table %s does not exist",  path[0]))
    }

    res, err := codec.MarshalJSONIndent(keeper.cdc, types.Table(table))
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

