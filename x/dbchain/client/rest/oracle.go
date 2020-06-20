package rest

import (
    "fmt"
    "time"
    "errors"
    "io/ioutil"
    "net/http"
    "encoding/json"
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

type IdCard struct {
    Name string     `json:"name"`
    IdNumber string `json:"id_number"`
}

const AliyunSmsKey    = "aliyun-sms-key"
const AliyunSmsSecret = "aliyun-sms-secret"
const AliyunSkyEyeAppCode = "aliyun-sky-eye-appcode"

var (
    associationMap = make(map[string]MobileVerfCode)
    aliyunSmsKey string
    aliyunSmsSecret string
)

func newIdCard(name, idNumber string) IdCard {
    return IdCard {
        Name:     name,
        IdNumber: idNumber,
    }
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
       
        //verificationCode := utils.GenerateVerfCode(6)
        verificationCode := "111111"
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
            saveToAuthTable(addr, "mobile", mobile)
            tryToSendToken(addr)
            rest.PostProcessResponse(w, cliCtx, "Success")
        } else {
            rest.WriteErrorResponse(w, http.StatusNotFound, "Failed to verify")
        }
    }
}

func oracleVerifyNameAndIdNumber(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        accessCode := vars["accessToken"]
        name       := vars["name"]
        idNumber   := vars["id_number"]

        addr, err := utils.VerifyAccessCode(accessCode)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }

        if VerifyNameAndIdNumber(name, idNumber) {
            saveToAuthTable(addr, "idcard", newIdCard(name, idNumber))
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
    return true
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

func saveToAuthTable(addr sdk.AccAddress, authType string, value interface{}) {
    rowFields := make(types.RowFields)
    rowFields["address"] = addr.String()
    rowFields["type"]    = authType
    bz, err := json.Marshal(value)
    if err != nil {
        fmt.Println("Failed to load oracle's account info!!!")
        return
    }
    rowFields["value"] = string(bz)
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

func VerifyNameAndIdNumber(name, id_number string) bool {
    url := fmt.Sprintf("https://personalbusiness.shumaidata.com/getPersonalBusiness?idcard=%s", id_number)
    request, err := http.NewRequest("GET", url, nil)
    if err != nil {
        fmt.Println("Failed to create new request!!!")
        return false
    }

    appCode := viper.GetString(AliyunSkyEyeAppCode)
    request.Header.Set("Authorization", "APPCODE " + appCode)

    var client = &http.Client{
        Timeout: time.Second * 10,
    }
    response, err := client.Do(request)
    if err != nil {
        fmt.Println("Failed to do request!!!")
        return false
    }

    defer response.Body.Close()
    if response.StatusCode != 200 {
        fmt.Printf("Returned code is %d!!!\n", response.StatusCode)
        return false
    }

    bodyBytes, _ := ioutil.ReadAll(response.Body)
    frName, err := getName(bodyBytes)    
    if err != nil {
        fmt.Printf("Returned code is %d!!!\n", response.StatusCode)
        return false
    }
    if name != frName {
        fmt.Printf("Names don't match")
        return false
    }
    return true
}

/////////////////////
//                 //
// helpe functions //
//                 //
/////////////////////

func unWrap(i interface{}, key string) interface{} {
   oneMap := i.(map[string]interface{})
   return oneMap[key]
}

func getName(bodyBytes []byte) (string, error) {
    var ent interface{}
    err := json.Unmarshal(bodyBytes, &ent)
    if err != nil {
        return "", errors.New("Failed to unmarshal body bytes")
    }
    var queue = []string{"data", "data", "ryposfrs"}
    for _, v := range queue {
        ent = unWrap(ent, v)
    }

    corps := ent.([]interface{})
    if len(corps) < 1 {
        return "", errors.New("No corps found")
    }
    corp := corps[0]
    c := corp.(map[string]interface{})
    name, ok := c["ryname"]
    if !ok {
        return "", errors.New("Couldn't find ryname")
    }
    return name.(string), nil
}
