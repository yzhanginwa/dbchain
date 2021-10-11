package oracle

import (
    "github.com/dbchaincloud/tendermint/crypto"
    "github.com/dbchaincloud/tendermint/crypto/algo"
    "github.com/dbchaincloud/tendermint/crypto/secp256k1"
    "github.com/spf13/viper"
    "github.com/mr-tron/base58"
    "github.com/dbchaincloud/tendermint/crypto/sm2"
)

const OracleEncryptedPrivKey = "oracle-encrypted-key"

var (
    oraclePrivKey crypto.PrivKey
    oraclePrivKeyLoaded = false
)
 
func LoadPrivKey() (crypto.PrivKey, error) {
    if oraclePrivKeyLoaded {
        return oraclePrivKey, nil
    }
    base58Str := viper.GetString(OracleEncryptedPrivKey)
    pkBytes, err:= base58.Decode(base58Str)
    if err != nil {
        return nil, err
    }
    switch algo.Algo {
    case algo.SM2:
        var privKey sm2.PrivKeySm2
        copy(privKey[:], pkBytes)
        oraclePrivKeyLoaded = true
        oraclePrivKey       = privKey
        return privKey, nil
    default:
        var privKey secp256k1.PrivKeySecp256k1
        copy(privKey[:], pkBytes)
        oraclePrivKeyLoaded = true
        oraclePrivKey       = privKey
        return privKey, nil
    }
}
