package keeper

import (
        "fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/yzhanginwa/rcv-chain/x/rcvchain/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the rcvservice Querier
const (
	QueryPoll    = "poll"
	QueryStatus  = "status"
	QueryTitles  = "titles"
	QueryBallot  = "ballot"
	QueryUserPolls = "userpolls"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryPoll:
			return queryPoll(ctx, path[1:], req, keeper)
		case QueryStatus:
			return queryStatus(ctx, path[1:], req, keeper)
		case QueryTitles:
			return queryTitles(ctx, []string{}, req, keeper)
		case QueryBallot:
			return queryBallot(ctx, path[1:], req, keeper)
		case QueryUserPolls:
			return queryUserPolls(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown nameservice query endpoint")
		}
	}
}

func queryPoll(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	poll, err := keeper.GetPoll(ctx, path[0])

	if err != nil {
		return []byte{}, sdk.ErrUnknownRequest(fmt.Sprintf("poll %s does not exist",  path[0]))
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, types.QueryResPoll(poll))
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

// nolint: unparam
func queryStatus(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	status := keeper.GetPollStatus(ctx, path[0])

	if status == "" {
		return []byte{}, sdk.ErrUnknownRequest(fmt.Sprintf("poll %s does not exist",  path[0]))
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, types.QueryResStatus{Status: status})
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

func queryTitles(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var titlesList types.QueryResTitles
	prefixLen := len(types.PollKeyPrefix)

        iterator := keeper.GetPollsIterator(ctx)

        for ; iterator.Valid(); iterator.Next() {
                k := string(iterator.Key())
		if len(k) > prefixLen {
			if k[:prefixLen] == types.PollKeyPrefix {
				id := k[prefixLen:]
				title := keeper.GetPollTitle(ctx, id)
				titlesList = append(titlesList, title)
			}
		}
        }

	res, err := codec.MarshalJSONIndent(keeper.cdc, titlesList)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

func queryBallot(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	if !keeper.IsPollPresent(ctx, path[0]) {
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("poll %s does not exist",  path[0]))
	}

	addr, err := sdk.AccAddressFromBech32(path[1])
	if err != nil {
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("address %s is invalid",  path[1]))
	}
		
	ballot, err := keeper.GetBallot(ctx, path[0], addr)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("%s",  err))
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, ballot)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

func queryUserPolls(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	addr, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("address %s is invalid",  path[0]))
	}

	userPoll, err := keeper.GetUserPolls(ctx, addr)
        if err != nil {
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("%s", err))
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, userPoll)
	if err != nil {
		panic("could not marshal resutl to JSON")
	}

	return res, nil
}
