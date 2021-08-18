package oracle

import (
	sdk "github.com/dbchaincloud/cosmos-sdk/types"
	"github.com/dbchaincloud/tendermint/crypto"
	"github.com/dbchaincloud/tendermint/crypto/sm2"
	"github.com/mr-tron/base58"
	"sync"
)
// this keys be used for these txs whose txHash need to be return
var oracleSpecialPkForNft = []string{
	"FpFFcuuuT3pEDxFUqLaTfigDo2oqhAE4zG2cmvdq5rxM",
	"6GVYcwR9g1F92JEKv3J4BstLgyRsoVKi7ihWeQb8RMSE",
	"63qLw9cb2sRwYgvXyVwQDbUcu2XEc5PzbiawquiTj3Yu",
	"4H9V8E3ApjjRZsLiTnZuMHh5zTporvCeTHdiNTzyAS65",
	"Abwk4R7aAYaQZSGP7EVKXJrQ9F1rqJuKzhtsbsgtrXGX",
	"23Lor6J2N1taweBiUMZS2JjGfEeRUgGx5g22o5UQ4Dab",
	"7TmRYJu7wwZXZt67MXoox9ehGaG9ToN7tgZGVwDdEWww",
	"FSSA47GYmUMZQGNmrMYeL34sidEzhuQEVkEt6JakUV6V",
	"EHtBbVvJXa2eRjeUtXt3WCM5w4smkvtuVfK7cu2cKhtF",
	"oAjj5RRfbgDvST9tZohXWwV24xxfw6NiodrZ7ijfAfm",
}

var count = 0
var mx sync.RWMutex
func loadSpecialPkForNtf() (crypto.PrivKey, sdk.AccAddress, error){
	mx.Lock()
	pkStr := oracleSpecialPkForNft[count]
	count = (count + 1) % len(oracleSpecialPkForNft)
	defer mx.Unlock()

	pkBytes, err:= base58.Decode(pkStr)
	if err != nil {
		return nil, nil, err
	}
	var privKey sm2.PrivKeySm2
	copy(privKey[:], pkBytes)
	addr := sdk.AccAddress(privKey.PubKey().Address())
	return privKey, addr, nil
}

