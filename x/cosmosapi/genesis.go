package cosmosapi

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type GenesisState struct {
	TablesRecords []Table  `json:"table_records"`
}

func NewGenesisState(tablesRecords []Table) GenesisState {
	return GenesisState{TablesRecords: tablesRecords}
}

func ValidateGenesis(data GenesisState) error {
	for _, record := range data.TablesRecords {
		if record.Owner == nil {
			return fmt.Errorf("invalid TablesRecord: Owner: %s. Error: Missing Owner", record.Owner)
		}
		if record.Name == "" {
			return fmt.Errorf("invalid TablesRecord: Name: %s. Error: Missing Value", record.Name)
		}
	}
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		TablesRecords: []Table{},
	}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) []abci.ValidatorUpdate {
	for _, record := range data.TablesRecords {
		keeper.CreateTable(ctx, record.Owner, record.Name, record.Fields)
	}
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	var records []Table
// TODO: update the following after implementing k.GetPollsIterator(ctx)
//	iterator := k.GetNamesIterator(ctx)
//	for ; iterator.Valid(); iterator.Next() {
//
//		name := string(iterator.Key())
//		whois := k.GetWhois(ctx, name)
//		records = append(records, whois)
//
//	}
	return GenesisState{TablesRecords: records}
}
