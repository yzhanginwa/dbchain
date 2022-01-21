package keeper

import (
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHandlePanic(t *testing.T) {
	//SafeHas
	describe := "SafeHas"
	has, err := tSafeHas(describe)
	require.Equal(t, false, has)
	require.Equal(t, errors.New(describe+" out of gas").Error(), err.Error())
	//SafeGet
	describe = "SafeGet"
	bz, err := tSafeGet(describe)
	require.Equal(t, []byte(nil), bz)
	require.Equal(t, errors.New(describe+" out of gas").Error(), err.Error())
	//SafeSet
	describe = "SafeSet"
	err = tSafeSet(describe)
	require.Equal(t, errors.New(describe+" out of gas").Error(), err.Error())
	//SafeIterator
	describe = "SafeIterator"
	it := tSafeIterator(describe)
	require.Equal(t, errors.New(describe+" out of gas").Error(), it.Error().Error())
}

func tSafeHas(describe string) (has bool, err error) {
	has = true
	err = errors.New("")
	defer handlePanic(describe, &err, &has)
	panicFunc()
	return has, err
}

func tSafeGet(describe string) (bz []byte, err error) {
	bz = []byte("hello")
	err = errors.New("")
	defer handlePanic(describe, &err, &bz)
	panicFunc()
	return bz, err
}

func tSafeSet(describe string) (err error) {
	err = errors.New("")
	defer handlePanic(describe, &err, nil)
	panicFunc()
	return err
}

func tSafeIterator(describe string) (it *SafeIterator) {
	it = NewSafeIterator(nil)
	defer handlePanic(describe, &it.err, nil)
	panicFunc()
	return it
}

func panicFunc() {
	panic(sdk.ErrorOutOfGas{"out of gas"})
}
