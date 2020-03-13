package utils

import (
    "testing"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/stretchr/testify/require"
)


func TestItemExists(t *testing.T) {
    cases := []struct {
        exist bool
        slice interface{}
        item  interface{}
    }{
         { true,  []string{"aaa", "bbb", "ccc"}, "aaa" },
         { false, []string{}, "aaa"},
         { true,  []uint{1, 2, 3}, uint(3)},
    }

    for _, tc := range cases {
        result := ItemExists(tc.slice, tc.item) 
        require.Equal(t, result, tc.exist)
    }
}


func TestAddressIncluded(t *testing.T) {
    acc1 := sdk.AccAddress([]byte("me"))
    acc2 := sdk.AccAddress([]byte("you"))

    cases := []struct {
        included bool
        slice    []sdk.AccAddress
        item     sdk.AccAddress
    }{
         { true,  []sdk.AccAddress{acc1, acc2}, acc1 },
         { false, []sdk.AccAddress{acc2}, acc1 },
    }

    for _, tc := range cases {
        result := AddressIncluded(tc.slice, tc.item) 
        require.Equal(t, result, tc.included)
    }
}

