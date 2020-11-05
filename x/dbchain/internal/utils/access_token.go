package utils

import (
    "fmt"
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

func MakeAccessCode(privKey secp256k1.PrivKeySecp256k1) string {
    now := time.Now().UnixNano() / 1000000
    timeStamp := strconv.Itoa(int(now))

    signature, err := privKey.Sign([]byte(timeStamp))
    if err != nil {
        panic("failed to sign timestamp")
    }

    pubKey := privKey.PubKey()
    pubKeyArray := pubKey.(secp256k1.PubKeySecp256k1)

    encodedPubKey := base58.Encode(pubKeyArray[:])
    encodedSig    := base58.Encode(signature)
    return fmt.Sprintf("%s:%s:%s", encodedPubKey, timeStamp, encodedSig)
}

func VerifyAccessCode(accessCode string) (sdk.AccAddress, error) {
    parts := strings.Split(accessCode, ":")
    if len(parts) != 3 {
        return nil, errors.New("Wrong access code format")
    }
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

