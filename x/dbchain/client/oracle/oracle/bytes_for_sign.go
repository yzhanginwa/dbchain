package oracle

import (
    "encoding/json"
    sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetSignBytes(msg UniversalMsg) []byte {
    jsonEncodedMsg, err := json.Marshal(msg)
    if err != nil {
        return []byte{}
    }
    return sdk.MustSortJSON([]byte(jsonEncodedMsg))
}

func StdSignBytes(chainID string, accnum uint64, sequence uint64, fee StdFee, msgs []UniversalMsg, memo string) []byte {
    msgsBytes := make([]json.RawMessage, 0, len(msgs))
    for _, msg := range msgs {
        msgsBytes = append(msgsBytes, json.RawMessage(GetSignBytes(msg)))
    }
    bz, err := json.Marshal(StdSignDoc{
        AccountNumber: accnum,
        ChainID:       chainID,
        Fee:           json.RawMessage(fee.Bytes()),
        Memo:          memo,
        Msgs:          msgsBytes,
        Sequence:      sequence,
    })
    if err != nil {
        panic(err)
    }
    return sdk.MustSortJSON(bz)
}
