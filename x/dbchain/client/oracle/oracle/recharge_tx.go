package oracle

import (
	"encoding/json"
	"fmt"
	"github.com/dbchaincloud/cosmos-sdk/client/context"
	sdk "github.com/dbchaincloud/cosmos-sdk/types"
	"github.com/dbchaincloud/tendermint/crypto"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	Undefined = 1
	Success   = 2
	Failed    = 3

)
func BuildAndSignBroadcastTx(cliCtx context.CLIContext, batch []UniversalMsg, privKey crypto.PrivKey, oracleAccAddr sdk.AccAddress) (string, int, string) {
	accNum, seq, err := GetAccountInfo(oracleAccAddr.String())
	if err != nil {
		fmt.Println("Failed to load oracle's account info!!!")
		return "", Failed, err.Error()
	}
	newBatch := make([]UniversalMsg, 0)
	for _, msg := range batch {
		if checkCanInsertRow(cliCtx, msg) {
			newBatch = append(newBatch, msg)
		}
	}

	txBytes, err := buildAndSignAndBuildTxBytes(cliCtx, newBatch, accNum, seq, privKey)
	if err != nil {
		return "", Failed, err.Error()
	}
	txHash := broadcastTxBytes(txBytes)

	var errInfo  = "undefined"
	var status int
	//waitUntilTxFinish(utils.MakeAccessCode(privKey), txHash)
	for i := 15; i > 0; i-- {
		//result := make(map[string]interface{})
		time.Sleep(1 * time.Second)
		status, errInfo = checkTxStatus(txHash)
		if status == Undefined {
			continue
		} else if status == Success {
			errInfo = ""
			break
		} else {
			break
		}
	}
	if status == Undefined || status == Failed {
		status = Failed
	}

	return txHash, status, errInfo
}

//status 1: query err
//       2: tx success
//       3: tx fail
func checkTxStatus(txHash string ) (int, string){
	resp, err := http.Get(fmt.Sprintf("http://localhost:1317/txs/%s", txHash))
	defer resp.Body.Close()
	if err != nil {
		return Undefined, err.Error()
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || body == nil {
		return Undefined, "undefined"
	}
	result := make(map[string]interface{})
	err = json.Unmarshal(body, &result)
	if err != nil {
		return Undefined, err.Error()
	}

	if info, ok := result["error"]; ok {
		return Undefined, info.(string)
	} else if rawLog, ok := result["raw_log"]; ok {
		log := rawLog.(string)
		data := make([]interface{}, 0)
		err := json.Unmarshal([]byte(log), &data)
		if err != nil {
			//tx fail
			// "insufficient funds: insufficient account funds; 18dbctoken \u003c 10000dbctoken: failed to execute message; message index: 0"
			return Failed, log
		} else {
			//tx success
			//[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"send"},{"key":"sender","value":"cosmos1n6yqmysvcz0cpnd52427ldjmjj493pk2uhpcu3"},{"key":"module","value":"bank"}]},{"type":"transfer","attributes":[{"key":"recipient","value":"cosmos156p5rmhpd3l709ygg7t80fu96fm4mrtsqhftvx"},{"key":"sender","value":"cosmos1n6yqmysvcz0cpnd52427ldjmjj493pk2uhpcu3"},{"key":"amount","value":"1dbctoken"}]}]}]
			return Success, ""
		}
	} else {
		return Undefined, "undefined"
	}

}