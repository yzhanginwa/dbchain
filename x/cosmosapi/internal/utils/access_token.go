package utils

import (
    "strings"
    "strconv"
    "errors"
    "time"
    "github.com/mr-tron/base58"
    "github.com/tendermint/tendermint/crypto/secp256k1"

    sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
    MaxAllowedTimeDiff = 15 * 1000   // 15 seconds
)


//////////////////
//              //
// helper funcs //
//              //
//////////////////

func VerifyAccessCode(accessCode string) (sdk.AccAddress, error) {
    parts := strings.Split(accessCode, ":")
    pubKeyBytes, _ := base58.Decode(parts[0])
    timeStamp      := parts[1]
    signature, _   := base58.Decode(parts[2])

    var pubKey secp256k1.PubKeySecp256k1
    copy(pubKey[:], pubKeyBytes)

    if ! pubKey.VerifyBytes([]byte(timeStamp), []byte(signature)) {
        return nil, errors.New("Failed to verify signature")
    }

    timeStampInt, err := strconv.Atoi(timeStamp)
    if err != nil {
        return nil, errors.New("Failed to verify access token")
    }
    now := time.Now().UnixNano() / 1000000
    diff := now - int64(timeStampInt)
    if diff < 0 { diff -= 0 }

    if diff < MaxAllowedTimeDiff {
        address := sdk.AccAddress(pubKey.Address())
        return address, nil
    } else {
        return nil, errors.New("Failed to verify access token")
    }
}

