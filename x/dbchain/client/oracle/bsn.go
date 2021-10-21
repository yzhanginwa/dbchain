package oracle

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cosmos/go-bip39"
	"github.com/dbchaincloud/cosmos-sdk/client/context"
	"github.com/dbchaincloud/cosmos-sdk/crypto/keys"
	sdk "github.com/dbchaincloud/cosmos-sdk/types"
	"github.com/dbchaincloud/tendermint/crypto"
	tmamino "github.com/dbchaincloud/tendermint/crypto/encoding/amino"
	"github.com/dbchaincloud/tendermint/crypto/sm2"
	"github.com/spf13/viper"
	"github.com/yzhanginwa/dbchain/x/dbchain/client/oracle/oracle"
	"github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
	"io/ioutil"
	"net/http"
	"strconv"
)

func applyAccountInfo(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		priv, secret, err := CreateMnemonic(keys.Sm2)
		if err != nil {
			generalResponse(w, map[string]string {"error " : "generate key pairs err"})
			return
		}
		pk := priv.(sm2.PrivKeySm2)
		pub := priv.PubKey().(sm2.PubKeySm2)

		add := sdk.AccAddress(pub.Address())

		data := map[string]string {
			"publicKey" : hex.EncodeToString(pub[:]),
			"privateKey" : hex.EncodeToString(pk[:]),
			"address" : add.String(),
			"mnemonic" : secret,
		}

		err = saveByOracle(cliCtx, data)
		if err != nil {
			generalResponse(w, map[string]string{ "error" : err.Error()})
			return
		}
		generalResponse(w, data)
		return

	}
}

func applyAccountInfoByPublicKey() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err  := ioutil.ReadAll(r.Body)
		if err != nil {
			generalResponse(w, map[string]string{"error" : err.Error()})
			return
		}
		postData := make(map[string]string)
		err = json.Unmarshal(result, &postData)
		if err != nil {
			generalResponse(w, map[string]string{"error" : err.Error()})
			return
		}
		publicKey := postData["publicKey"]
		pubBytes , err := hex.DecodeString(publicKey)
		if err != nil {
			generalResponse(w, map[string]string{"error" : err.Error()})
			return
		}

		if len(pubBytes) != 33 {
			generalResponse(w, map[string]string{"error" : "publicKey length error"})
			return
		}
		var pubKey sm2.PubKeySm2
		copy(pubKey[:], pubBytes)

		add := sdk.AccAddress(pubKey.Address())
		data := map[string]string {
			"publicKey" : publicKey,
			"address" : add.String(),
		}
		generalResponse(w, data)
		return

	}
}

func rechargeTx(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := readBodyData(r)
		if err != nil {
			generalResponse(w, map[string]string{ "error" : err.Error()})
		}
		bsnAddress := data["bsnAddress"]
		userAccountAddress := data["userAccountAddress"]
		rechargeGas := data["rechargeGas"]
		tx, status, errInfo := sendFromBsnAddressToUserAddress(cliCtx, storeName, bsnAddress, userAccountAddress, rechargeGas)
		generalResponse(w, map[string]interface{}{
			"txHash" : tx,
			"state" : status,
			"remarks" : errInfo,
		})
		return
	}
}

func getAccountTxByTimeOrByHeight(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := readBodyData(r)
		if err != nil {
			generalResponse(w, map[string]string{ "error" : err.Error()})
		}

		userAccountAddress := data["userAccountAddress"]
		startDate := data["startDate"]
		endDate := data["endDate"]
		if startDate == "" || endDate == "" {
			startDate, endDate = getQueryDate(cliCtx, data)
			if startDate == "" || endDate == "" {
				generalResponse(w, map[string]string{"error" : "parameters err"})
				return
			}
		}
		if userAccountAddress == "" || startDate == "" || endDate == "" {
			generalResponse(w, map[string]string{"error" : "expect 3 parameters : userAccountAddress, startDate, endDate"})
			return
		}

		year, month, day, hour, minite, second := 0,0,0,0,0,0
		nStartDate , _ := fmt.Sscanf(startDate,"%d-%d-%d %d:%d:%d", &year, &month, &day, hour, minite, second)
		nEndDate , _ := fmt.Sscanf(endDate,"%d-%d-%d %d:%d:%d", &year, &month, &day)
		if nStartDate != 3 || nEndDate != 3 {
			generalResponse(w, map[string]string{"error" : "time format error, it should be  yyyy-mm-dd"})
			return
		}



		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/account_txs_by_time/%s/%s/%s", storeName, userAccountAddress, startDate, endDate), nil)
		if err != nil {
			generalResponse(w, map[string]string{"error" : err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
		return
	}
}

func getQueryDate(cliCtx context.CLIContext, data map[string]string) (string, string) {
	startHeight := data["startBlockHeight"]
	endHeight := data["endBlockHeight"]
	if startHeight == "" || endHeight == "" {
		return "", ""
	}
	node , err := cliCtx.GetNode()
	if err != nil {
		return "", ""
	}
	start , err := strconv.ParseInt(startHeight,10, 64)
	if err != nil {
		return "", ""
	}
	end , err := strconv.ParseInt(endHeight,10, 64)
	if err != nil {
		return "", ""
	}
	startBlock,err := node.Block(&start)
	if err != nil {
		return "", ""
	}
	endBlock , err := node.Block(&end)
	if err != nil {
		return "", ""
	}
	startTime := startBlock.Block.Time.Local().Format("2006-01-02 15:04:05")
	endTime := endBlock.Block.Time.Local().Format("2006-01-02 15:04:05")
	return startTime, endTime
}
///////////////////
//               //
//   help func   //
//               //
///////////////////

//aes key

var secret []byte
var oraclePrivateKey crypto.PrivKey
const secretKey = "secret_key"

func generalResponse(w http.ResponseWriter, data interface{}) {
	bz,_ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(bz)
}

func CreateMnemonic(algo keys.SigningAlgo) (crypto.PrivKey, string, error) {

	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return  nil, "", err
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, "", err
	}

	//
	derivedPriv, err := keys.StdDeriveKey(mnemonic, keys.DefaultBIP39Passphrase, sdk.GetConfig().GetFullFundraiserPath(), algo)
	if err != nil {
		return nil, "", err
	}

	privKey, err := keys.StdPrivKeyGen(derivedPriv, algo)
	if err != nil {
		return nil, "", err
	}

	return privKey, mnemonic, nil
}

func readBodyData(r *http.Request) (map[string]string, error) {
	result, err  := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil , err
	}

	postData := make(map[string]string)
	err = json.Unmarshal(result, &postData)
	if err != nil {
		return nil, err
	}
	return postData, nil
}

func sendFromBsnAddressToUserAddress(cliCtx context.CLIContext, storeName, bsnAddress, userAccountAddress, rechargeGas string) (string, int, string){
	from, err  := sdk.AccAddressFromBech32(bsnAddress)
	if err != nil {
		return "", oracle.Failed, err.Error()
	}
	to, err := sdk.AccAddressFromBech32(userAccountAddress)
	if err != nil {
		return "", oracle.Failed, err.Error()
	}
	coins, err := sdk.ParseCoins(rechargeGas)
	if err != nil {
		return "", oracle.Failed, err.Error()
	}

	pk , err := loadUserPrivateKeyFromChain(cliCtx, storeName, bsnAddress)
	if err != nil {
		if err != nil {
			return "", oracle.Failed, err.Error()
		}
	}

	msg := oracle.NewMsgSend(from, to, coins)
	txHash, status, errInfo := oracle.BuildAndSignBroadcastTx(cliCtx, []oracle.UniversalMsg{msg}, pk, from)
	fmt.Println(txHash, status, errInfo)
	return txHash, status, errInfo
}

func saveByOracle( cliCtx context.CLIContext, data map[string]string ) error {

	pk , err := loadOraclePrivateKey()
	if err != nil {
		return err
	}
	aes, err := loadAesEncryptKey()
	if err != nil {
		return err
	}

	user := data["address"]
	bz , err := json.Marshal(map[string]string {
		"address" : data["address"],
		"privateKey" : data["privateKey"],
	})
	if err != nil {
		return err
	}

	ecryptBz ,err := AESEncrypt(bz, aes)
	if err != nil {
		return err
	}
	hexBz := hex.EncodeToString(ecryptBz)

	addr := sdk.AccAddress(pk.PubKey().Address())


	msg := types.NewMsgSaveUserPrivateKey(addr, user, hexBz)
	err = msg.ValidateBasic()
	if err != nil {
		return err
	}
	oracle.BuildTxsAndBroadcast(cliCtx,  []oracle.UniversalMsg{msg})
	return nil
}

func loadOraclePrivateKey() (crypto.PrivKey,error) {

	if oraclePrivateKey == nil {
		var err  error
		oraclePrivateKey , err  = oracle.LoadPrivKey()
		if err != nil {
			return nil, err
		}
	}

	return oraclePrivateKey, nil
}

func loadAesEncryptKey() ([]byte, error) {

	if secret == nil {
		key := viper.GetString(secretKey)
		if key == "" {
			return nil, errors.New("secretKey is empty")
		}
		return []byte(key), nil
	}
	return secret, nil
}

func loadUserPrivateKeyFromChain(cliCtx context.CLIContext, storeName, addr string) (crypto.PrivKey, error) {

	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/get_user_private_key/%s", storeName, addr), nil)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, errors.New("from address does not exit")
	}
	aes, err := loadAesEncryptKey()
	if err != nil {
		return nil, err
	}

	bz, _ := hex.DecodeString(string(res))
	data , err := AESDecrypt(bz, aes)
	if err != nil {
		return nil, err
	}
	keyInfo := make(map[string]string, 0)
	err = json.Unmarshal(data, &keyInfo)
	if err != nil {
		return nil, err
	}
	pkStr := keyInfo["privateKey"]
	pkBytes , _ := hex.DecodeString(pkStr)
	private, _  := tmamino.PrivKeyFromBytes(pkBytes)
	return private, nil
}