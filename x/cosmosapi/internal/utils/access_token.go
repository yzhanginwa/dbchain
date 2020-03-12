package utils

import (
    "strings"
    "strconv"
    "errors"
    "time"
    "regexp"
    //"encoding/hex"
    "encoding/base64"
    "github.com/tendermint/tendermint/crypto/secp256k1"

    sdk "github.com/cosmos/cosmos-sdk/types"
    //"github.com/yzhanginwa/cosmos-api/x/cosmosapi/internal/types"
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
    r1 := regexp.MustCompile("-")
    r2 := regexp.MustCompile("_")
    accessCode1 := r1.ReplaceAllString(accessCode, "+");
    accessCode2 := r2.ReplaceAllString(accessCode1, "/");

    parts := strings.Split(accessCode2, ":")
    pubKeyBytes, _ := base64.StdEncoding.DecodeString(parts[0])
    timeStamp      := parts[1]
    signature, _   := base64.StdEncoding.DecodeString(parts[2])

    //pubKeyBytes, _ := hex.DecodeString(pubKeyStr)
    //pubKey, _ := crypto.PubKey(hex.DecodeString(pubKeyStr))

    var pubKey secp256k1.PubKeySecp256k1
    copy(pubKey[:], pubKeyBytes)
    //pubKey := crypto.PubKey(pubKeyBytes)

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

