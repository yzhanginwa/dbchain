package utils

import (
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/mr-tron/base58"
)

const OracleEncryptedPrivKey = "oracle-encrypted-key"

var (
	oraclePrivKey secp256k1.PrivKeySecp256k1
	oraclePrivKeyLoaded = false
)

func LoadPrivKey() (secp256k1.PrivKeySecp256k1, error) {
	if oraclePrivKeyLoaded {
		return oraclePrivKey, nil
	}
	base58Str := viper.GetString(OracleEncryptedPrivKey)
	pkBytes, err:= base58.Decode(base58Str)
	if err != nil {
		return secp256k1.PrivKeySecp256k1{}, err
	}
	var privKey secp256k1.PrivKeySecp256k1
	copy(privKey[:], pkBytes)
	oraclePrivKeyLoaded = true
	oraclePrivKey       = privKey
	return privKey, nil
}
