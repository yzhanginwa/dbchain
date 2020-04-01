package keeper 

import (
    "fmt"
    "testing"
    "regexp"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/stretchr/testify/require"
)

func TestGenerateAppCode(t *testing.T) {
    var pattern = regexp.MustCompile("^[A-Z0-9]{10}$")

    addr1 := sdk.AccAddress([]byte("short one"))
    addr2 := sdk.AccAddress([]byte("this is a long string"))
    addr3 := sdk.AccAddress([]byte("this is a another long string"))

    cases := []struct {
        matching bool
        address  sdk.AccAddress
    }{
         { true,  addr1},
         { true,  addr2 },
         { true,  addr3 },
    }

    for _, tc := range cases {
        code := generateNewAppCode(tc.address)
        fmt.Println(code) 
        result := pattern.MatchString(code)
        require.Equal(t, result, tc.matching)
    }
}

