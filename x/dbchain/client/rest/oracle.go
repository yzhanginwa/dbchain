package rest

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "github.com/gorilla/mux"
    "github.com/spf13/viper"
    "github.com/mr-tron/base58"
    "github.com/tendermint/tendermint/crypto/secp256k1"
    "github.com/cosmos/cosmos-sdk/client/context"
    "github.com/cosmos/cosmos-sdk/types/rest"
    sdk "github.com/cosmos/cosmos-sdk/types"
    authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
    "github.com/cosmos/cosmos-sdk/x/auth/exported"
    amino "github.com/tendermint/go-amino"
    cryptoamino "github.com/tendermint/tendermint/crypto/encoding/amino"

    rpcclient "github.com/tendermint/tendermint/rpc/client"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"

    "github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

type MobileVerfCode struct {
    Mobile string     `json:"mobile"`
    VerfCode string   `json:"verf_code"`
}

const OracleEncryptedPrivKey = "oracle-encrypted-key"
const AliyunSmsKey    = "aliyun-sms-key"
const AliyunSmsSecret = "aliyun-sms-secret"

var (
    aminoCdc = amino.NewCodec()
    associationMap = make(map[string]MobileVerfCode)
    oraclePrivKey secp256k1.PrivKeySecp256k1
    oraclePrivKeyLoaded = false
    aliyunSmsKey string
    aliyunSmsSecret string
)

func init () {
    aminoCdc.RegisterInterface((*sdk.Msg)(nil), nil)
    aminoCdc.RegisterInterface((*sdk.Tx)(nil), nil)
    aminoCdc.RegisterConcrete(types.MsgInsertRow{}, "dbchain/InsertRow", nil)
    cryptoamino.RegisterAmino(aminoCdc)
    authtypes.RegisterCodec(aminoCdc)
}

func oracleSendVerfCode(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        accessCode := vars["accessToken"]
        mobile     := vars["mobile"]

        addr, err := utils.VerifyAccessCode(accessCode)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }
       
        verificationCode := utils.GenerateVerfCode(6)
        cacheMobileAndVerificationCode(addr.String(), mobile, verificationCode)
        if sent := sendVerificationCode(mobile, verificationCode); !sent {
            rest.WriteErrorResponse(w, http.StatusNotFound, "Failed to send sms")
            return
        } 
        rest.PostProcessResponse(w, cliCtx, "Success")
    }
}

func oracleVerifyVerfCode(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        accessCode      := vars["accessToken"]
        mobile          := vars["mobile"]
        verificationCode := vars["verificationCode"]

        addr, err := utils.VerifyAccessCode(accessCode)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }

        if VerifyVerfCode(addr.String(), mobile, verificationCode) {
            saveToAuthTable(addr, mobile)
            rest.PostProcessResponse(w, cliCtx, "Success")
        } else {
            rest.WriteErrorResponse(w, http.StatusNotFound, "Failed to verify")
        }
    }
}

//////////////////////
//                  //
// helper functions //
//                  //
//////////////////////

func LoadAliyunSmsKeyAndSecret() (string, string) {
    key := viper.GetString(AliyunSmsKey)
    secret := viper.GetString(AliyunSmsSecret)
    return key, secret
}
 
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

func cacheMobileAndVerificationCode(strAddr string, mobile string, verificationCode string) bool {
    mobileVerfCode := MobileVerfCode {
        Mobile: mobile,
        VerfCode: verificationCode,
    }

    associationMap[strAddr] = mobileVerfCode
    return true
}   

func sendVerificationCode(mobile string, verificationCode string) bool {
    if aliyunSmsKey == "" {
        aliyunSmsKey, aliyunSmsSecret = LoadAliyunSmsKeyAndSecret()
    }

    // aliyun sms service
    client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", aliyunSmsKey, aliyunSmsSecret)
    request := dysmsapi.CreateSendSmsRequest()
    request.Scheme = "https"

    request.SignName = "大中华区块链公章"
    request.TemplateCode = "SMS_192960014"
    request.PhoneNumbers = mobile
    request.TemplateParam = fmt.Sprintf("{\"code\": \"%s\"}", verificationCode)

    response, err := client.SendSms(request)
    if err != nil {
        fmt.Print(err.Error())
        return false
    }
    fmt.Printf("response is %#v\n", response)

    //to send verification through a provider
    return true
}

func VerifyVerfCode(strAddr string , mobile string, verificationCode string) bool {
    if mobileCode, ok := associationMap[strAddr]; ok {
        delete(associationMap, strAddr)
        if mobileCode.Mobile == mobile && mobileCode.VerfCode == verificationCode {
            return true
        }
    }
    return false
}

// this function is called in oracle. so it needs to broadcast a msg to save data to blockchain
func saveToAuthTable(addr sdk.AccAddress, mobile string) {
    privKey, err := LoadPrivKey()
    if err != nil {
        fmt.Println("Failed to load oracle's private key!!!")
        return
    }

    oracleAccAddr := sdk.AccAddress(privKey.PubKey().Address())

    accNum, seq, err := getAccountInfo(oracleAccAddr.String())
    if err != nil {
        fmt.Println("Failed to load oracle's account info!!!")
        return
    }
   
    fmt.Printf("\nAccount number: %d, seq: %d\n\n", accNum, seq)

    rowFields := make(types.RowFields)
    rowFields["address"] = addr.String()
    rowFields["type"]    = "mobile"
    rowFields["value"]   = mobile

    rowFieldsJson, err := json.Marshal(rowFields)
    if err != nil { 
        fmt.Println("Oracle: Failed to to json.Marshal!!!")
        return 
    }
    
    msg := types.NewMsgInsertRow(sdk.AccAddress(privKey.PubKey().Address()), "0000000001", "authentication", rowFieldsJson)
    err = msg.ValidateBasic()
    if err != nil {
        fmt.Println("Oracle: Failed validate new message!!!")
        return
    }

    msgs := []sdk.Msg{msg}
    stdFee := authtypes.NewStdFee(200000, sdk.Coins{sdk.NewCoin("dbctoken", sdk.NewInt(1))})

    stdSignMsg := authtypes.StdSignMsg{
        ChainID:       "testnet",
        AccountNumber: accNum,
        Sequence:      seq,
        Memo:          "",
        Msgs:          msgs,
        Fee:           stdFee,
    }

    sig, err := privKey.Sign(stdSignMsg.Bytes())
    if err != nil {
        fmt.Println("Oracle: Failed to sign message!!!")
        return
    }

    stdSignature :=authtypes.StdSignature{
        PubKey:    privKey.PubKey(),
        Signature: sig,
    }

    newStdTx := authtypes.NewStdTx(msgs, stdFee, []authtypes.StdSignature{stdSignature}, "")

    encoder := authtypes.DefaultTxEncoder(aminoCdc)
    txBytes, err := encoder(newStdTx)
    if err != nil {
        fmt.Println("Oracle: Failed to marshal StdTx!!!")
        return
    }

    //cliCtx.BroadcastTxAsync(txBytes)
    rpc, err := rpcclient.NewHTTP("http://localhost:26657", "/websocket")
    if err != nil {
        fmt.Printf("failted to get client: %v\n", err)
        return
    }

    _, err = rpc.BroadcastTxAsync(txBytes)
    if err != nil {
        fmt.Printf("failted to broadcast transaction: %v\n", err)
        return
    }
}

func getAccountInfo(address string) (uint64, uint64, error) {
    resp, err := http.Get(fmt.Sprintf("http://localhost:1317/auth/accounts/%s", address))
    if err != nil {
        fmt.Println("failed to get account info") 
        return 0, 0, err
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)

    type MyAccount struct {
      Height string            `json:"height"`   
      Result exported.Account  `json:"result"`
    }

    var account MyAccount
    if err := aminoCdc.UnmarshalJSON(body, &account); err != nil {
        fmt.Printf("failted to broadcast unmarshal account body\n")
        return 0, 0, err
    }

    seq := account.Result.GetSequence()
    accountNumber := account.Result.GetAccountNumber()

    return accountNumber, seq, nil
}
