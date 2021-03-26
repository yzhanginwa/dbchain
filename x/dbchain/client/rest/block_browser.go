package rest

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"sync"
	"time"
)

var (
	exp *regexp.Regexp
	txsStatistic *txStatistic
	totalTxs  *totalStatistic
)

const (
	oneDaySeconds = 86400
	blockOutTime = 5
)

func init(){
	regExpression := `[0-9]{4}-[0-9]{2}-[0-9]{2}`
	exp , _ = regexp.Compile(regExpression)
	txsStatistic = NewTxStatistic()
	totalTxs = NewTotalStatistic()
}

func showCurrentDayTxsNum(cliCtx context.CLIContext) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		txs, _ := getOneDayTxs(cliCtx, 0, false)
		bz , _ := json.Marshal(txs)
		rest.PostProcessResponse(w, cliCtx, bz)
	}
}

func showRecentDaysTxsNum(cliCtx context.CLIContext) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
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
					tx, date := getOneDayTxs(cliCtx,i, true)
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
	}
}

func showTotalTxsNum(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//excluding data of the day
		txs := 0
		daysAgo := 1
		updateTimeStamp := getCurrentDayStartTimeStamp()
		//TODO concurrent query
		//if we want to concurrent query, we need now how many days need to query
		for {
			tx, timeStamp := getOneDayTxs(cliCtx, daysAgo, false)
			daysAgo++
			if totalTxs.compare(timeStamp) {
				txs += tx
			} else if timeStamp == 0 && totalTxs.isZreo(){
				txs += tx
				totalTxs.update(updateTimeStamp,txs)
				break
			} else {
				totalTxs.update(updateTimeStamp,txs)
				break
			}
		}

		bz , _ := json.Marshal(totalTxs.getTotalTxs())
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
func getOneDayTxs(cliCtx context.CLIContext, daysAgo int, upDateCache bool) (int,int64) {
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
	length         int
}

func NewTxStatistic() *txStatistic {
	var TxStatistic txStatistic
	TxStatistic.data = make(map[int64]int)
	TxStatistic.length = cacheDays
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
	if len(T.data)  <= T.length {
		return
	}
	days := T.sortDateBigToSmall()
	expirationData := days[0]
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