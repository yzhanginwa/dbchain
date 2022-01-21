package keeper

import (
	"errors"
	"github.com/cosmos/cosmos-sdk/store/gaskv"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"reflect"
)

//this file is wrap of gasKv store for handling panic of out of gas

type (
	Iterator = sdk.Iterator
	KVStore  = sdk.KVStore
)

type SafeStore struct {
	KVStore
}

func DbChainStore(ctx sdk.Context,storeKey sdk.StoreKey) *SafeStore{
	rawStore := ctx.KVStore(storeKey)
	return NewSafeStore(rawStore)
}

func DbChainStoreWithOutGas(ctx sdk.Context,storeKey sdk.StoreKey) *SafeStore{
	gasConfig := sdk.GasConfig {
		HasCost:          0,
		DeleteCost:       0,
		ReadCostFlat:     0,
		ReadCostPerByte:  0,
		WriteCostFlat:    0,
		WriteCostPerByte: 0,
		IterNextCostFlat: 0,
	}
	rawStore := gaskv.NewStore(ctx.MultiStore().GetKVStore(storeKey), ctx.GasMeter(), gasConfig)
	return NewSafeStore(rawStore)
}

func NewSafeStore(store KVStore) *SafeStore {
	return &SafeStore{store}
}

func (S *SafeStore) Get(key []byte) (bz []byte, err error) {

	defer handlePanic("SafeGet", &err, &bz)
	bz = S.KVStore.Get(key)
	return bz, nil

}

func (S *SafeStore) Delete(key []byte) (err error) {

	defer handlePanic("SafeDelete", &err, nil)
	S.KVStore.Delete(key)
	return nil
}

func (S *SafeStore) Set(key []byte, value []byte) (err error) {

	defer handlePanic("SafeSet", &err, nil)
	S.KVStore.Set(key, value)
	return nil
}

func (S *SafeStore) Has(key []byte) (has bool, err error) {

	defer handlePanic("SafeHas", &err, &has)
	has = S.KVStore.Has(key)
	return has, nil
}

func (S *SafeStore)Iterator(start, end []byte) *SafeIterator{
	iter := S.KVStore.Iterator(start,end)
	return NewSafeIterator(iter)
}
//Wrapping an iterator

type SafeIterator struct {
	err error
	Iterator
}

func NewSafeIterator(it sdk.Iterator) *SafeIterator {
	return &SafeIterator{
		err:      nil,
		Iterator: it,
	}
}

func (sIter *SafeIterator) Next() {

	defer handlePanic("SafeIterator", &sIter.err, nil)
	sIter.Iterator.Next()
}

func (sIter *SafeIterator) Error() error {
	return sIter.err
}

func handlePanic(describe string, err *error, data interface{}) {
	r := recover()
	if r == nil {
		return
	}
	//set err
	switch r.(type) {
	case sdk.ErrorOutOfGas:
		*err = errors.New(describe + " out of gas")
	case sdk.ErrorGasOverflow:
		*err = errors.New(describe + " ErrorGasOverflow")
	default:
		*err = errors.New(describe + " SafeHas default err")
	}
	//set value
	value := reflect.ValueOf(data)
	for {
		if value.Kind() != reflect.Ptr {
			break
		}
		value = value.Elem()
	}

	switch value.Kind() {
	case reflect.Bool:
		value.SetBool(false)
	case reflect.Slice:
		value.SetBytes(nil)
	case reflect.Int:
		value.SetInt(0)
	default:
		//do nothing
	}

}
