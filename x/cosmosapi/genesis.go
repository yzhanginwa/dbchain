package rcvchain

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type GenesisState struct {
	PollRecords []Poll  `json:"poll_records"`
}

func NewGenesisState(pollRecords []Poll) GenesisState {
	return GenesisState{PollRecords: pollRecords}
}

func ValidateGenesis(data GenesisState) error {
	for _, record := range data.PollRecords {
		if record.Owner == nil {
			return fmt.Errorf("invalid PollRecord: Title: %s. Error: Missing Owner", record.Title)
		}
		if record.Title == "" {
			return fmt.Errorf("invalid PollRecord: Owner: %s. Error: Missing Value", record.Owner)
		}
	}
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		PollRecords: []Poll{},
	}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) []abci.ValidatorUpdate {
	for _, record := range data.PollRecords {
		keeper.CreatePoll(ctx, record.Title, record.Owner)
	}
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	var records []Poll
// TODO: update the following after implementing k.GetPollsIterator(ctx)
//	iterator := k.GetNamesIterator(ctx)
//	for ; iterator.Valid(); iterator.Next() {
//
//		name := string(iterator.Key())
//		whois := k.GetWhois(ctx, name)
//		records = append(records, whois)
//
//	}
	return GenesisState{PollRecords: records}
}
