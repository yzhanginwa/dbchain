package dbchain

import (
	"encoding/hex"
	"errors"
	sdk "github.com/dbchaincloud/cosmos-sdk/types"
	stdtypes "github.com/dbchaincloud/cosmos-sdk/x/auth/types"
	"github.com/dbchaincloud/tendermint/crypto/tmhash"
	"github.com/yzhanginwa/dbchain/x/dbchain/internal/keeper/cache"
	"github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
	"time"
)



func setTxStatus(ctx sdk.Context, keeper Keeper, err *error)  {
	txHash := hex.EncodeToString(tmhash.Sum(ctx.TxBytes()))
	cdc := keeper.GetCdc()
	decoder := stdtypes.DefaultTxDecoder(cdc)
	tx,decoderErr := decoder(ctx.TxBytes())
	if decoderErr != nil {
		Err := errors.New("decode tx err")
		timeStamp := time.Now().Unix()
		tempStatus := types.NewTxStatus("fail",0, Err.Error(),timeStamp)
		cache.TxStatusCache.Store(txHash,tempStatus)
		return
	}
	msgs := tx.GetMsgs()
	currentStatesIm,ok := cache.TxStatusCache.Load(txHash)
	//there may be multi msgs in one tx
	if ok {
		currentStatus := currentStatesIm.(*types.TxStatus)
		currentStatus.Index++
		if *err != nil{
			currentStatus.State = cache.TxStateFail
			currentStatus.Err = (*err).Error()
		} else if currentStatus.Index == len(msgs) {
			currentStatus.State = cache.TxStateSuccess
		} else {
			currentStatus.State = cache.TxStateProcessing
		}
	} else {
		timeStamp := time.Now().Unix()
		tempStatus := types.NewTxStatus("",1, "",timeStamp)
		if *err != nil {
			tempStatus.State = "fail"
			tempStatus.Err = (*err).Error()
		} else if len(msgs) == 1 {
			tempStatus.State = cache.TxStateSuccess
		} else {
			tempStatus.State = cache.TxStateProcessing
		}
		cache.TxStatusCache.Store(txHash, tempStatus)
	}
	return
}


func TxCacheInvalid() {
	tk := time.NewTicker(time.Second * cache.TxStateInvalidTime)
	for {
		select {
		case <-tk.C:
			cacheInvalid()
		}
	}
}

func cacheInvalid() {
	tk := time.NewTicker(cache.TxInvalidCheckRunTime)
	nowTime := time.Now().Unix()
	cache.TxStatusCache.Range(func(key, value interface{}) bool {
			select {
			case <-tk.C://time out
				return false
			default:
				data := value.(*types.TxStatus)
				if nowTime - data.GetTimeStamp() > cache.TxStateInvalidTime  {
					cache.TxStatusCache.Delete(key)
				}
			}
			return true
	})
}
