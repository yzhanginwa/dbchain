package utils

import (
	"reflect"
        "bytes"
	sdk "github.com/dbchaincloud/cosmos-sdk/types"
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

func StringIncluded(strSlice []string, str string) bool {
    for _, item := range strSlice {
        if item == str {
            return true
        }
    }
    return false
}

func RemoveStringFromSet(set []string, item string) []string {
    for i, v := range set {
        if v == item {
            set[i] = set[len(set)-1]
            set[len(set)-1] = ""    // probably keep from mem leaking
            return set[:len(set)-1]
        }
    }
    return set
}
