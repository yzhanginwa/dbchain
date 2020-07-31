package types

import (
    "testing"

    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/stretchr/testify/require"
)

func TestMsgCreateTable(t *testing.T) {
    owner := sdk.AccAddress([]byte("me"))
    msg := NewMsgCreateTable(owner, "appCode", "table_foo", []string{"fld1", "fld2"})

    require.Equal(t, msg.Route(), RouterKey)
    require.Equal(t, msg.Type(), "create_table")
}

func TestMsgCreateTableValidation(t *testing.T) {
    owner := sdk.AccAddress([]byte("me"))

    cases := []struct {
        valid bool
        tx    MsgCreateTable
    }{
        {true, NewMsgCreateTable(owner, "appCode", "table_foo", []string{"fld1"})},
        {true, NewMsgCreateTable(owner, "appCode", "table_foo", []string{"fld1", "fld2"})},
        {false, NewMsgCreateTable(owner, "appCode", "", []string{"fld1", "fld2"})},
        {false, NewMsgCreateTable(owner, "appCode", "table_foo", []string{})},
    }

    for _, tc := range cases {
        err := tc.tx.ValidateBasic()
        if tc.valid {
            require.Nil(t, err)
        } else {
            require.NotNil(t, err)
        }
    }
}

func TestMsgCreateTableGetSignBytes(t *testing.T) {
    owner := sdk.AccAddress([]byte("me"))
    msg := NewMsgCreateTable(owner, "ABCDEFGH", "foo", []string{"fld1", "fld2"})
    res := msg.GetSignBytes()

    expected := `{"type":"dbchain/CreateTable","value":{"app_code":"ABCDEFGH","fields":["fld1","fld2"],"owner":"cosmos1d4js690r9j","table_name":"foo"}}`

    require.Equal(t, expected, string(res))
}
