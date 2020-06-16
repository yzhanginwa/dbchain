package rest

import (
    "fmt"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/spf13/viper"
    "github.com/cosmos/cosmos-sdk/client/context"
    "github.com/cosmos/cosmos-sdk/types/rest"
    sdk "github.com/cosmos/cosmos-sdk/types"

    "github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/client/rest/oracle"

    "github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

type MobileVerfCode struct {
    Mobile string     `json:"mobile"`
    VerfCode string   `json:"verf_code"`
}

const AliyunSmsKey    = "aliyun-sms-key"
const AliyunSmsSecret = "aliyun-sms-secret"

var (
    associationMap = make(map[string]MobileVerfCode)
    aliyunSmsKey string
    aliyunSmsSecret string
)

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
            tryToSendToken(addr)
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
    //TODO: use logger
    fmt.Printf("response is %#v\n", response)
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

func saveToAuthTable(addr sdk.AccAddress, mobile string) {
    rowFields := make(types.RowFields)
    rowFields["address"] = addr.String()
    rowFields["type"]    = "mobile"
    rowFields["value"]   = mobile

    oracle.InsertRow("0000000001", "authentication", rowFields)
}

func tryToSendToken(addr sdk.AccAddress) {
    accNum, _, err := oracle.GetAccountInfo(addr.String())
    if err != nil {
        fmt.Println("Failed to load oracle's account info!!!")
        return
    }

    if accNum == 0 {
        oracle.SendFirstTokenTo(addr)
    }
}
