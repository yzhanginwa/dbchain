package oracle

import (
	"encoding/hex"
	"encoding/json"
	"github.com/dbchaincloud/cosmos-sdk/client/context"
	sdk "github.com/dbchaincloud/cosmos-sdk/types"
	"github.com/dbchaincloud/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/yzhanginwa/dbchain/x/dbchain/client/oracle/oracle"
	"github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"sync"
	"time"
)

var (
	exp *regexp.Regexp
)

const (
	oneDaySeconds = 86400
	blockOutTime = 5
)

func init(){
	regExpression := `[0-9]{4}-[0-9]{2}-[0-9]{2}`
	exp, _ = regexp.Compile(regExpression)
}

func showCurrentDayTxsNum(cliCtx context.CLIContext) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		TxsStatistic := loadTxStatistic(cliCtx)
		txs, _ := getOneDayTxs(cliCtx, 0, false, TxsStatistic)
		bz , _ := json.Marshal(txs)
		rest.PostProcessResponse(w, cliCtx, bz)
	}
}

func showRecentDaysTxsNum(cliCtx context.CLIContext) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		TxsStatistic := loadTxStatistic(cliCtx)

		TxsStatisticCopy := make(map[int64]int)
		for k,v := range TxsStatistic.data {
			TxsStatisticCopy[k] = v
		}

		vars := mux.Vars(r)
		strDaysAgo := vars["daysAgo"]
		daysAgo, err  := strconv.Atoi(strDaysAgo)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "daysAgo parameter")
			return
		}
		if daysAgo <= 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "daysAgo should not be 0")
			return
		}
		if daysAgo >= 15 {
			daysAgo = 15
		}

		dates := make(chan int64, daysAgo)
		txs   := make(chan int, daysAgo)
		go func() {
			for i := 1; i <= daysAgo; i++ {
				go func(i int) {
					tx, date := getOneDayTxs(cliCtx,i, true, TxsStatistic)
					dates <- date
					txs <- tx
				}(i)
			}
		}()

		result := make(map[int64]int)
		loopCount := daysAgo
		for date := range dates {
			txs := <-txs
			result[date] = txs
			loopCount--
			if loopCount <= 0 {
				break
			}
		}
		close(dates)
		close(txs)
		bz , _ := json.Marshal(result)
		rest.PostProcessResponse(w, cliCtx, bz)
		//Check whether it needs to be stored in the database
		for k, v := range TxsStatistic.data {
			if TxsStatisticCopy[k] != v {
				endProcessing(cliCtx, TxsStatistic)
				break
			}
		}
	}
}

func showTotalTxsNum(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//excluding data of the day

		TotalTxs := loadTotalStatistic(cliCtx)
		TxsStatistic  := loadTxStatistic(cliCtx)

		var TotalTxsCopy = NewTotalStatistic()
		TotalTxsCopy.date = TotalTxs.date
		TotalTxsCopy.txNum = TotalTxs.txNum

		txs := 0
		daysAgo := 1
		updateTimeStamp := getCurrentDayStartTimeStamp()

		if TotalTxs.date >= updateTimeStamp {
			bz , _ := json.Marshal(TotalTxs.getTotalTxs())
			rest.PostProcessResponse(w, cliCtx, bz)
			return
		}
		for {
			tx, timeStamp := getOneDayTxs(cliCtx, daysAgo, false, TxsStatistic)
			daysAgo++
			if TotalTxs.date < timeStamp {
				txs += tx
			} else if TotalTxs.date == timeStamp{
				txs += tx
				TotalTxs.update(updateTimeStamp,txs)
				break
			}else if timeStamp == 0 && TotalTxs.isZreo(){
				txs += tx
				TotalTxs.update(updateTimeStamp,txs)
				break
			} else {
				TotalTxs.update(updateTimeStamp,txs)
				break
			}
		}

		bz , _ := json.Marshal(TotalTxs.getTotalTxs())
		rest.PostProcessResponse(w, cliCtx, bz)
		//Check whether it needs to be stored in the database
		if TotalTxsCopy.date != TotalTxs.date {
			endProcessing(cliCtx ,TotalTxs)
		}
	}
}

func showBlockTxsHash(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		height, err  := strconv.ParseInt(vars["height"],10,64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "block height err ")
			return
		}
		node, err := cliCtx.GetNode()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "GetNode err : " + err.Error())
			return
		}
		block, err := node.Block(&height)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "get block err : " + err.Error())
			return
		}
		result := make([]string,0)
		Txs := block.Block.Txs
		for _,tx := range Txs {
			txha := hex.EncodeToString(tx.Hash())
			result = append(result, txha)
		}
		bz , _ := json.Marshal(result)
		rest.PostProcessResponse(w, cliCtx, bz)
	}
}

func showAllAccounts(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := cliCtx.QueryWithData("custom/dbchain/allAccounts", nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func showAllApplications(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := cliCtx.QueryWithData("custom/dbchain/application_browser", nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

//From the current time to a certain day
func getOneDayTxs(cliCtx context.CLIContext, daysAgo int, upDateCache bool, txsStatistic *txStatistic) (int,int64) {
	txs := 0
	//get from cache first
	startTimeStamp := getStartTimeStamp(daysAgo)
	txs = txsStatistic.GetOneDayTx(startTimeStamp)
	if txs == -1 {
		txs = 0
	} else {
		return txs, startTimeStamp
	}
	//
	node, err := cliCtx.GetNode()
	if err != nil {
		return -1, startTimeStamp
	}
	currentBlock, err := node.Block(nil)
	if err != nil {
		return -1, startTimeStamp
	}
	currentHeight := currentBlock.Block.Height

	startHeight , endHeight := getStartAndEndBlockHeight(startTimeStamp, currentHeight)
	if endHeight == 0 {
		return 0, 0
	}

	numGO := make(chan bool, 30)
	txChan := make(chan int, 10)
	go func() {
		for i := startHeight; i <= endHeight; i++ {
			var height int64 = i
			numGO <- true
			go func(i int64) {
				block , err := node.BlockResults(&height)
				if err != nil {
					<-numGO
					txChan <- 0
					return
				}
				txChan <- len(block.TxsResults)
				<-numGO
			}(i)
		}
	}()

	var loopCount int64 = endHeight - startHeight + 1
	for tx := range txChan {
		txs += tx
		loopCount--
		if loopCount <= 0 {
			break
		}
	}
	close(numGO)
	close(txChan)
	if upDateCache {
		txsStatistic.Update(startTimeStamp,txs)
	}
	return txs, startTimeStamp
}

func getStartAndEndBlockHeight(startTimeStamp, currentHeight int64) (startHeight, endHeight int64){

	startHeight = currentHeight - (time.Now().Unix() - startTimeStamp) / blockOutTime
	endHeight = startHeight + oneDaySeconds/blockOutTime
	if startHeight < 0 {
		startHeight = 0
	}
	if endHeight > currentHeight {
		endHeight = currentHeight
	} else if endHeight < 0 {
		endHeight = 0
	}
	return startHeight, endHeight
}

func getStartTimeStamp(daysAgo int) int64 {
	CurrentDayStartTimeStamp := getCurrentDayStartTimeStamp()
	if daysAgo == 0 {
		return CurrentDayStartTimeStamp
	}
	startTimeStamp := CurrentDayStartTimeStamp - int64(daysAgo) * oneDaySeconds
	return startTimeStamp
}

////////////////////////////////////
//                                //
//                                //
//           help func            //
//                                //
//                                //
////////////////////////////////////
func getCurrentDayStartTimeStamp() int64 {
	nowTime := time.Now()
	format := nowTime.Format(time.RFC3339)
	suffix := `T00:00:00+08:00`
	prefix := exp.FindString(format)
	newTime, _  := time.Parse(time.RFC3339, prefix + suffix)
	return  newTime.Unix()
}

const (
	cacheDays = 15	//Statistic days
)

//////////////////////////////////////
//                                  //
//           txStatistic            //
//                                  //
//////////////////////////////////////


type txStatistic struct {
	rwMutex sync.RWMutex
	data map[int64]int
}

func NewTxStatistic() *txStatistic {
	var TxStatistic txStatistic
	TxStatistic.data = make(map[int64]int)
	return &TxStatistic
}

func (T *txStatistic)GetOneDayTx(date int64) int{
	T.rwMutex.RLock()
	defer T.rwMutex.RUnlock()
	tx , ok := T.data[date]
	if !ok {
		return -1
	}
	return tx
}

func (T *txStatistic)Update(date int64, txs int) {
	T.rwMutex.Lock()
	defer T.rwMutex.Unlock()
	T.data[date] = txs
	if len(T.data)  <= cacheDays {
		return
	}

	expirationData := T.getOldestDay()
	delete(T.data, expirationData)
	return
}

func (T *txStatistic) sortDateBigToSmall() []int64{
	stores := make([]int64,0)
	for k ,_ := range T.data {
		stores = append(stores, k)
	}
	//Sort from big to small
	sort.Slice(stores, func(i, j int) bool {
		return  stores[i] > stores[j]
	})
	return stores
}

func (T *txStatistic) getOldestDay() int64 {
	var day int64
	for dayTime , _ := range T.data {
			if day == 0 || day > dayTime {
					day = dayTime
				}
		}
	return day
}

//////////////////////////////////////
//                                  //
//           total txs              //
//                                  //
//////////////////////////////////////
type totalStatistic struct {
	txNum int
	date int64
}

func NewTotalStatistic() *totalStatistic {
	TotalStatistic := totalStatistic{
		txNum: 0,
		date: 0,
	}
	return &TotalStatistic
}

func (total *totalStatistic)update(date int64, txs int){
	total.txNum += txs
	total.date = date
	return
}

func (total *totalStatistic)compare(date int64) bool {
	if total.date >= date {
		return false
	}
	return true
}

func (total *totalStatistic)isZreo()bool {
	return total.date == 0
}

func (total *totalStatistic)getTotalTxs()int{
	return total.txNum
}


func loadTotalStatistic(cliCtx context.CLIContext) *totalStatistic{
	TotalTxs := NewTotalStatistic()
	out := queryTotalTxs(cliCtx)
	if out == nil {
		return TotalTxs
	}
	TotalTxs.txNum = int(out["txNum"])
	TotalTxs.date  = out["date"]
	return TotalTxs
}

func loadTxStatistic(cliCtx context.CLIContext) *txStatistic{
	var TxStatistic = NewTxStatistic()

	out := queryTxStatistic(cliCtx)
	if out == nil {
		return TxStatistic
	}
	TxStatistic.data = out
	return TxStatistic
}



func endProcessing(cliCtx context.CLIContext, data interface{}) {
	priv, _  := oracle.LoadPrivKey()
	var oracleAddr = sdk.AccAddress(priv.PubKey().Address())
	msgs := make([]oracle.UniversalMsg, 0)

	switch data := data.(type) {
	case *totalStatistic:
		msgs = append(msgs, types.NewMsgUpdateTotalTx(oracleAddr, marshalTotalTxs(data)))
	case *txStatistic:
		msgs = append(msgs, types.NewMsgUpdateTxStatistic(oracleAddr, marshalTxsStatistic(data)))
	default:
		return
	}
	oracle.BuildTxsAndBroadcast(cliCtx, msgs)
}

///////////////////////////////
//                           //
//        help func          //
//                           //
///////////////////////////////

func marshalTotalTxs(totalTxs *totalStatistic) string {
	temp := make(map[string]int64)
	temp["txNum"] = int64(totalTxs.txNum)
	temp["date"] = totalTxs.date
	bz , _ := json.Marshal(temp)
	return string(bz)
}

func marshalTxsStatistic(txsStatistic *txStatistic) string {
	bz, _ := json.Marshal(txsStatistic.data)
	return string(bz)
}

func queryTxStatistic(cliCtx context.CLIContext)  map[int64]int {
	res, _, err := cliCtx.QueryWithData("custom/dbchain/dbchainRecentTxNum", nil)
	if err != nil {
		return nil
	}
	var out map[int64]int
	err = json.Unmarshal(res, &out)
	if err != nil || out == nil {
		return nil
	}
	return out
}

func queryTotalTxs(cliCtx context.CLIContext)  map[string]int64 {
	res, _, err := cliCtx.QueryWithData("custom/dbchain/dbchainTxNum", nil)
	if err != nil {
		return nil
	}
	var out map[string]int64
	err = json.Unmarshal(res, &out)
	if err != nil || out == nil {
		return nil
	}
	return out
}