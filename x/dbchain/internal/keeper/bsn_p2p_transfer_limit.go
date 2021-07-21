package keeper

import (
	"encoding/json"
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"

)


func (k Keeper) ShowChainSuperAdmins(ctx sdk.Context, addr sdk.Address) []string {
	admins := k.getChainSuperAdmins(ctx, addr)
	// only admin can query all admins
	for _, admin := range admins {
		if addr.String() == admin {
			return admins
		}
	}
	return nil
}

func (k Keeper) getChainSuperAdmins(ctx sdk.Context, addr sdk.Address) []string {
	admins := make([]string, 0)
	store := DbChainStore(ctx, k.storeKey)
	key := getChainSuperAdminsKey()
	bz, err := store.Get([]byte(key))
	if err != nil || bz == nil {
		return admins
	}
	err = json.Unmarshal(bz, &admins)
	if err != nil {
		return nil
	}
	return admins
}

func (k Keeper)  IsChainSuperAdmin(ctx sdk.Context, addr sdk.Address) bool {
	admins := k.ShowChainSuperAdmins(ctx, addr)
	if len(admins) > 0 {
		return true
	}
	return false
}

func (k Keeper) ModifyMemberOfAdmins(ctx sdk.Context, modifier , addr sdk.Address, action string) error {

	admins := k.getChainSuperAdmins(ctx, modifier)
	if len(admins) == 0 {
		admins = make([]string, 0)
	} else if !k.IsChainSuperAdmin(ctx, modifier) {
		return errors.New("permission forbidden")
	}

	if action == "add" {
		for _, admin := range admins {
			if admin == addr.String() {
				return errors.New("address already exists")
			}
		}
		admins = append(admins, addr.String())
	} else {
		adminExist := false
		for i, admin := range admins {
			if admin == addr.String() {
				adminExist = true
				admins = append(admins[:i], admins[i+1:]...)
				break
			}
		}
		if !adminExist {
			return errors.New("address dose not exist")
		}
	}
	bz , err := json.Marshal(admins)
	if err != nil {
		return err
	}
	store := DbChainStore(ctx, k.storeKey)
	key := getChainSuperAdminsKey()
	return store.Set([]byte(key), bz)
}

func (k Keeper) SetP2PTransferLimit(ctx sdk.Context, modifier sdk.Address, limit bool) error {
	if !k.IsChainSuperAdmin(ctx, modifier) {
		return errors.New("permission forbidden")
	}
	store := DbChainStore(ctx, k.storeKey)
	key := getP2PTransferLimit()
	bz, err := store.Get([]byte(key))
	if err != nil {
		return err
	}

	if bz == nil {
		data , err := json.Marshal(limit)
		if err != nil {
			return err
		}
		return store.Set([]byte(key), data)
	}

	var current bool
	err = json.Unmarshal(bz, & current)
	if err != nil {
		return err
	}

	if current == limit {
		info := fmt.Sprintf("current status already is %v", limit)
		return errors.New(info)
	}

	data , err := json.Marshal(limit)
	if err != nil {
		return err
	}
	return store.Set([]byte(key), data)
}

func (k Keeper)ShowCurrentLimitP2PTransferStatus(ctx sdk.Context) bool {

	store := DbChainStore(ctx, k.storeKey)
	key := getP2PTransferLimit()
	bz, err := store.Get([]byte(key))
	if err != nil || bz == nil {
		return false
	}

	var current bool
	err = json.Unmarshal(bz, & current)
	if err != nil {
		return false
	}
	return current
}
