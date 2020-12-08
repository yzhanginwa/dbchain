package oracle

import (
    "github.com/spf13/viper"
    "github.com/mr-tron/base58"
    "github.com/dbchaincloud/tendermint/crypto/sm2"
)

const OracleEncryptedPrivKey = "oracle-encrypted-key"

var (
    oraclePrivKey sm2.PrivKeySm2
    oraclePrivKeyLoaded = false
)
 
func LoadPrivKey() (sm2.PrivKeySm2, error) {
    if oraclePrivKeyLoaded {
        return oraclePrivKey, nil
    }
    base58Str := viper.GetString(OracleEncryptedPrivKey)
    pkBytes, err:= base58.Decode(base58Str)
    if err != nil {
        return sm2.PrivKeySm2{}, err
    }
    var privKey sm2.PrivKeySm2
    copy(privKey[:], pkBytes)
    oraclePrivKeyLoaded = true
    oraclePrivKey       = privKey
    return privKey, nil
}
