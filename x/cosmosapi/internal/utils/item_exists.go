package utils

import (
	"reflect"
        "bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ItemExists(slice interface{}, item interface{}) bool {
    s := reflect.ValueOf(slice)

    if s.Kind() != reflect.Slice {
        panic("Invalid data-type")
    }

    for i := 0; i < s.Len(); i++ {
        if s.Index(i).Interface() == item {
            return true
        }
    }

    return false
}

func AddressIncluded(addresses []sdk.AccAddress, address sdk.AccAddress) bool {
    for _, addr := range addresses {
        if bytes.Compare(address, addr) == 0 {
            return true
        }
    }
    return false
}
