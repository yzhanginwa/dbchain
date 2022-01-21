package utils

import (
    "testing"
    "bytes"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/stretchr/testify/require"
    "github.com/tendermint/tendermint/crypto/secp256k1"
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

func TestStringIncluded(t *testing.T) {
    cases := []struct {
        included bool
        slice []string
        item  string
    }{
         { true,  []string{"aaa", "bbb", "ccc"}, "aaa" },
         { false, []string{"aaa"}, "ccc"},
         { false, []string{}, "aaa"},
    }

    for _, tc := range cases {
        result := StringIncluded(tc.slice, tc.item)
        require.Equal(t, result, tc.included)
    }
}

func TestSplitFieldName(t *testing.T) {
    tableName, ok := GetTableNameFromForeignKey("supplier_id")
    require.Equal(t, ok, true)
    require.Equal(t, tableName, "supplier")

    tableName, ok = GetTableNameFromForeignKey("_id")
    require.Equal(t, ok, false)
}

func TestConvertIntToByteArray(t *testing.T) {
    cases := []struct {
        matching bool
        number int64
        byteArray []byte
    }{
         { false, 1,   []byte{1, 0, 0, 0, 0, 0, 0, 0} },
         { true,  1,   []byte{0, 0, 0, 0, 0, 0, 0, 1} },
         { true,  255, []byte{0, 0, 0, 0, 0, 0, 0, 255} },
         { true,  256, []byte{0, 0, 0, 0, 0, 0, 1, 0} },
         { true,  257, []byte{0, 0, 0, 0, 0, 0, 1, 1} },
    }

    for _, tc := range cases {
        ary := IntToByteArray(tc.number)
        result := (bytes.Compare(ary, tc.byteArray) == 0)
        require.Equal(t, tc.matching, result)

        n := ByteArrayToInt(tc.byteArray)
        require.Equal(t, tc.matching, (n == tc.number))
    }
}

func TestAccessToken(t *testing.T) {
    privKey := secp256k1.GenPrivKey()
    accAddr := sdk.AccAddress(privKey.PubKey().Address())

    accessCode := MakeAccessCode(privKey)

    addr, err := VerifyAccessCode(accessCode)
    require.Equal(t, true, (err == nil))
    require.Equal(t, true, (accAddr.String() == addr.String()))
}
