package keeper

import (
	sdk "github.com/dbchaincloud/cosmos-sdk/types"
)

func (k Keeper) UpdateTotalTx(ctx sdk.Context, data string)  error {
	store := DbChainStore(ctx, k.storeKey)

	keyAt := getTotalTx()
	store.Set([]byte(keyAt), k.cdc.MustMarshalBinaryBare(data))
	return  nil
}

func (k Keeper) UpdateTxStatistic(ctx sdk.Context, data string) error{
	store := DbChainStore(ctx, k.storeKey)

	keyAt := getTxStatistic()
	store.Set([]byte(keyAt), k.cdc.MustMarshalBinaryBare(data))
	return nil
}


func (k Keeper) GetDbchainTxNum(ctx sdk.Context) ([]byte, error){
	store := DbChainStore(ctx, k.storeKey)
	keyAt := getTotalTx()
	bz, err := store.Get([]byte(keyAt))
	if err != nil {
		return nil, err
	}
	if bz == nil {
		return []byte{}, nil
	}
	var data  string
	k.cdc.MustUnmarshalBinaryBare(bz, &data)
	return []byte(data),nil
}

func (k Keeper) GetDbchainRecentTxNum(ctx sdk.Context) ([]byte, error){
	store := DbChainStore(ctx, k.storeKey)
	keyAt := getTxStatistic()
	bz, err := store.Get([]byte(keyAt))
	if err != nil {
		return nil, err
	}
	if bz == nil {
		return []byte{}, nil
	}
	var data  string
	k.cdc.MustUnmarshalBinaryBare(bz, &data)
	return []byte(data),nil
}