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
	"FpFFcuuuT3pEDxFUqLaTfigDo2oqhAE4zG2cmvdq5rxM",  //cosmos1g3mngn3r5s70wq6ztc32qlpjz07rmg6k95t08k
	"6GVYcwR9g1F92JEKv3J4BstLgyRsoVKi7ihWeQb8RMSE",  //cosmos1p23zpga40qmw7swmkl5pu80r6upfz8qg8wyjk9
	"63qLw9cb2sRwYgvXyVwQDbUcu2XEc5PzbiawquiTj3Yu",  //cosmos1x268y3m4z6pec9lmds3g4l4lxy3twlvfj39k5h
	"4H9V8E3ApjjRZsLiTnZuMHh5zTporvCeTHdiNTzyAS65",  //cosmos102c2e8jlx5d8a38e9evgurp9vun60ygtlp5vvq
	"Abwk4R7aAYaQZSGP7EVKXJrQ9F1rqJuKzhtsbsgtrXGX",  //cosmos1vp930yhsy2y2q4harcqkddnd42lnwq9qkzw9qy
	"23Lor6J2N1taweBiUMZS2JjGfEeRUgGx5g22o5UQ4Dab",  //cosmos1sr9c3wk5uyqcdzv9upzzpfmn4gw5djn8ly0sag
	"7TmRYJu7wwZXZt67MXoox9ehGaG9ToN7tgZGVwDdEWww",  //cosmos17h8vtgq3cvarxqpwqfd30yhrp9agg2rxrhlv3v
	"FSSA47GYmUMZQGNmrMYeL34sidEzhuQEVkEt6JakUV6V",  //cosmos1pntfyuyj3p8myx5hjp25egrcv7xpjlwpmpfkt8
	"EHtBbVvJXa2eRjeUtXt3WCM5w4smkvtuVfK7cu2cKhtF",  //cosmos1xw2hywsmstywmvrchrt7lukzkv5y6r7ywmg25s
	"oAjj5RRfbgDvST9tZohXWwV24xxfw6NiodrZ7ijfAfm",   //cosmos1x54wkthwf9jwcu2j7sm6p08x66vhachnjs9zp7
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

