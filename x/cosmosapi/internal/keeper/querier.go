package keeper

import (
    "fmt"
    "github.com/cosmos/cosmos-sdk/codec"

    sdk "github.com/cosmos/cosmos-sdk/types"
    abci "github.com/tendermint/tendermint/abci/types"
    //"github.com/yzhanginwa/cosmos-api/x/cosmosapi/internal/types"
)

// query endpoints supported by the cosmosapi service Querier
const (
    QueryTables   = "tables"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
    return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
        switch path[0] {
        case QueryTables:
            if len(path) > 1 {
                return queryTable(ctx, path[1:], req, keeper)
            } else {
                return queryTables(ctx, req, keeper)
            }
        default:
            return nil, sdk.ErrUnknownRequest("unknown cosmosapi query endpoint")
        }
    }
}

func queryTables(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
    tables, err := keeper.getTables(ctx)

    if err != nil {
        return []byte{}, sdk.ErrUnknownRequest("Can not get table names")
    }

    res, err := codec.MarshalJSONIndent(keeper.cdc, tables)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryTable(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
    table, err := keeper.GetTable(ctx, path[0])

    if err != nil {
        return []byte{}, sdk.ErrUnknownRequest(fmt.Sprintf("Table %s does not exist",  path[0]))
    }

    res, err := codec.MarshalJSONIndent(keeper.cdc, table)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

