package oracle

import (
    "fmt"
    "github.com/cosmos/cosmos-sdk/client/context"
    "github.com/gorilla/mux"
    //"github.com/yzhanginwa/dbchain/x/dbchain/client/oracle"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
    r.HandleFunc(fmt.Sprintf("/%s/upload/{%s}", storeName, "accessToken"), uploadFileHandler(cliCtx)).Methods("POST")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/send_verf_code/{%s}/{%s}", storeName, "accessToken", "mobile"), oracleSendVerfCode(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/verify_verf_code/{%s}/{%s}/{%s}", storeName, "accessToken", "mobile", "verificationCode"), oracleVerifyVerfCode(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/verify_name_and_id_number/{%s}/{%s}/{%s}", storeName, "accessToken", "name", "id_number"), oracleVerifyNameAndIdNumber(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/verify_corp_info/{%s}/{%s}/{%s}/{%s}", storeName, "accessToken", "corp_name", "reg_number", "credit_code"), oracleVerifyCorpInfo(cliCtx, storeName)).Methods("GET")

    r.HandleFunc(fmt.Sprintf("/%s/oracle/new_app_user/{%s}", storeName, "accessToken"), appNewOneCoin(cliCtx, storeName)).Methods("GET")

    //dbpay
    r.HandleFunc(fmt.Sprintf("/%s/oracle/dbcpay/{%s}/recipient_address", storeName, "accessToken"), oracleRecipientAddress(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/dbcpay/{%s}/{%s}", storeName, "accessToken", "payType"), oracleCallDbcPay(cliCtx, storeName)).Methods("POST")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/dbcpay_query/{%s}", storeName, "accessToken"), oracleQueryPayStatus(cliCtx, storeName)).Methods("POST")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/dbcpay_notify", storeName), oracleSavePayStatus(cliCtx, storeName)).Methods("POST")


    //block browser. do not need access token
    r.HandleFunc(fmt.Sprintf("/%s/oracle/blockchain/txs_num/current_day", storeName), showCurrentDayTxsNum(cliCtx)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/blockchain/txs_num/recent_day/{daysAgo}", storeName), showRecentDaysTxsNum(cliCtx)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/blockchain/txs_num/total", storeName), showTotalTxsNum(cliCtx)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/blockchain/all_accounts", storeName), showAllAccounts(cliCtx)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/blockchain/all_applications", storeName), showAllApplications(cliCtx)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/block/txs_hash/{%s}", storeName, "height"), showBlockTxsHash(cliCtx)).Methods("GET")

}
