package oracle

import (
    "fmt"
    "github.com/dbchaincloud/cosmos-sdk/client/context"
    "github.com/gorilla/mux"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
    r.HandleFunc(fmt.Sprintf("/%s/upload/{%s}/{%s}", storeName, "accessToken", "appCode"), uploadFileHandler(cliCtx)).Methods("POST")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/nft/send_verf_code/{%s}/{%s}", storeName, "purpose", "mobile"), oracleSendVerfCode(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/verify_verf_code/{%s}/{%s}", storeName, "mobile", "verificationCode"), oracleVerifyVerfCode(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/verify_name_and_id_number/{%s}/{%s}/{%s}", storeName, "accessToken", "name", "id_number"), oracleVerifyNameAndIdNumber(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/verify_corp_info/{%s}/{%s}/{%s}/{%s}", storeName, "accessToken", "corp_name", "reg_number", "credit_code"), oracleVerifyCorpInfo(cliCtx, storeName)).Methods("GET")

    r.HandleFunc(fmt.Sprintf("/%s/oracle/new_app_user/{%s}", storeName, "accessToken"), appNewOneCoin(cliCtx, storeName)).Methods("GET")

    //dbpay
    r.HandleFunc(fmt.Sprintf("/%s/oracle/dbcpay/{%s}/recipient_address", storeName, "accessToken"), oracleRecipientAddress(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/dbcpay/{%s}/{%s}", storeName, "accessToken", "payType"), oracleCallDbcPay(cliCtx, storeName)).Methods("POST")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/dbcpay_query/{%s}", storeName, "accessToken"), oracleQueryPayStatus(cliCtx, storeName)).Methods("POST")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/dbcpay_notify", storeName), oracleSavePayStatus(cliCtx, storeName)).Methods("POST")
    //payment
    r.HandleFunc(fmt.Sprintf("/%s/oracle/payment/dbcpay/{%s}/recipient_address", storeName, "accessToken"), oracleRecipientAddress(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/payment/{%s}/{%s}", storeName, "accessToken", "payType"), oracleCallDbcPay(cliCtx, storeName)).Methods("POST")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/applepay/{%s}", storeName, "accessToken"), oracleApplepay(cliCtx, storeName)).Methods("POST")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/payment_query/{%s}", storeName, "accessToken"), oracleQueryPayStatus(cliCtx, storeName)).Methods("POST")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/payment_notify", storeName), oracleSavePayStatus(cliCtx, storeName)).Methods("POST")

    //block browser. do not need access token
    r.HandleFunc(fmt.Sprintf("/%s/oracle/blockchain/txs_num/current_day", storeName), showCurrentDayTxsNum(cliCtx)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/blockchain/txs_num/recent_day/{daysAgo}", storeName), showRecentDaysTxsNum(cliCtx)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/blockchain/txs_num/total", storeName), showTotalTxsNum(cliCtx)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/blockchain/all_accounts", storeName), showAllAccounts(cliCtx)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/blockchain/all_applications", storeName), showAllApplications(cliCtx)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/block/txs_hash/{%s}", storeName, "height"), showBlockTxsHash(cliCtx)).Methods("GET")

    //
    r.HandleFunc(fmt.Sprintf("/%s/oracle/organization/get_secrect_key", storeName), organizationGetSecretKey(cliCtx)).Methods("POST")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/user/get_secrect_key", storeName), userGetSecretKey(cliCtx)).Methods("POST")
    // organization verify . param : hash_key  state
    r.HandleFunc(fmt.Sprintf("/%s/oracle/organization/verify", storeName), organizationVerify(cliCtx)).Methods("POST")
    //user verify
    r.HandleFunc(fmt.Sprintf("/%s/oracle/user/verify", storeName), userVerifyCode(cliCtx,false)).Methods("POST")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/user/comform", storeName), userVerifyCode(cliCtx, true)).Methods("POST")
    //user destroy
    r.HandleFunc(fmt.Sprintf("/%s/oracle/user/destroy", storeName), userDestoryCode(cliCtx)).Methods("POST")

    // bsn open interface
    //密钥托管创建链账户
    r.HandleFunc(fmt.Sprintf("/%s/oracle/bsn/account/apply", storeName), applyAccountInfo(cliCtx)).Methods("POST")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/bsn/account/apply/publicKey", storeName), applyAccountInfoByPublicKey()).Methods("POST")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/bsn/account/recharge", storeName), rechargeTx(cliCtx, storeName)).Methods("POST")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/bsn/account/tx", storeName), getAccountTxByTimeOrByHeight(cliCtx, storeName)).Methods("POST")

    //nft
    //register
    r.HandleFunc(fmt.Sprintf("/%s/oracle/nft/register", storeName), nftUserRegister(cliCtx, storeName)).Methods("POST")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/nft/login", storeName), nftUserLogin(cliCtx, storeName)).Methods("POST")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/nft/reset_password", storeName), nftUserResetPassword(cliCtx, storeName)).Methods("POST")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/nft/real_name_authentication", storeName), realNameAuthentication(cliCtx, storeName)).Methods("POST")

    r.HandleFunc(fmt.Sprintf("/%s/oracle/nft/nft_make_before", storeName), nftMakeBefore(cliCtx, storeName)).Methods("POST")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/nft/nft_make_without_money", storeName), nftMakeOld(cliCtx, storeName)).Methods("POST")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/nft/nft_public", storeName), nftPublish(cliCtx, storeName)).Methods("POST")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/nft/nft_withdraw", storeName), nftWithdraw(cliCtx, storeName)).Methods("POST")
    //only published can be purchased
    r.HandleFunc(fmt.Sprintf("/%s/oracle/nft/nft_buy", storeName), nftBuy(cliCtx, storeName)).Methods("POST")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/nft/nft_save_receipt", storeName), nftSaveReceipt(cliCtx, storeName)).Methods("POST")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/nft/nft_save_receipt_initiative", storeName), nftSaveReceiptInitiative(cliCtx, storeName)).Methods("POST")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/nft/nft_transfer", storeName), nftTransfer(cliCtx, storeName)).Methods("POST")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/nft/nft_edit_personnal_information", storeName), nftEditPersonalInformation(cliCtx, storeName)).Methods("POST")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/nft/nft_collect", storeName), nftCollect(cliCtx, storeName)).Methods("POST")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/nft/nft_cancle_collect", storeName), nftCancelCollect(cliCtx, storeName)).Methods("POST")
    //query
    r.HandleFunc(fmt.Sprintf("/%s/oracle/nft/find/{%s}/{%s}", storeName, "name", "id"), nftFindById(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/nft/find_by/{%s}/{%s}/{%s}", storeName, "name", "field", "value"), nftFindByField(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/nft/find_all/{%s}", storeName, "name"), nftFindAll(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/nft/querier/{%s}", storeName, "querierBase58"), nftFindByQuerier(cliCtx, storeName)).Methods("GET")

    r.HandleFunc(fmt.Sprintf("/%s/oracle/nft/popular_author/{%s}", storeName, "numbers"), nftFindPopularAuthor(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/nft/lastest_nft/{%s}/{%s}", storeName, "page","numbers"), nftFindLastestNft(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/nft/nft_details/{%s}", storeName, "denom_id"), nftFindNftDetails(cliCtx, storeName)).Methods("GET")
    //get data from session
    r.HandleFunc(fmt.Sprintf("/%s/oracle/nft/user_info", storeName), nftUserInfo(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/nft/user_income/{%s}", storeName, "user_id"), nftUserIncome(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/nft/user_collect/{%s}/{%s}", storeName, "user_id", "number_or_detail"), nftUserCollectInfo(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/nft/user/all/token_record", storeName), nftUserAllTokenRecord(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/nft/user/invite/token_record", storeName), nftUserInvitationRecord(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/nft/user_nft_number/{%s}", storeName, "user_id"), nftUserNftNumber(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/nft/user_nft_order_number", storeName), nftUserNftOrderNumber(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/nft/nfts_of_user_make/{%s}/{%s}", storeName, "publish_status","user_id"), nftsOfUserMake(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/nft/nfts_of_user_buy", storeName), nftsOfUserBuy(cliCtx, storeName)).Methods("GET")
    r.HandleFunc(fmt.Sprintf("/%s/oracle/nft/user_basic_info/{%s}", storeName, "user_id"), nftOfUserBasicInfo(cliCtx, storeName)).Methods("GET")
}
