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

type Mobile struct {
    Mobile string   `json:"mobile"`
}

type IdCard struct {
    Name string     `json:"name"`
    IdNumber string `json:"id_number"`
}

type CorpInfo struct {
    CorpName string   `json:"corp_name"`
    RegNumber string  `json:"reg_number"`
    CreditCode string `json:"credit_code"`
}

const AliyunSmsKey    = "aliyun-sms-key"
const AliyunSmsSecret = "aliyun-sms-secret"
const AliyunSkyEyeAppCode = "aliyun-sky-eye-appcode"

var (
    associationMap = make(map[string]MobileVerfCode)
    aliyunSmsKey string
    aliyunSmsSecret string
)

func newMobile(mobile string) Mobile {
    return Mobile {
        Mobile: mobile,
    }
}

func newIdCard(name, idNumber string) IdCard {
    return IdCard {
        Name:     name,
        IdNumber: idNumber,
    }
}

func newCorpInfo(corpName, regNumber, creditCode string) CorpInfo {
    return CorpInfo {
        CorpName:   corpName,
        RegNumber:  regNumber,
        CreditCode: creditCode,
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
            saveToAuthTable(addr, "mobile", newMobile(mobile))
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

        if verifyNameAndIdNumber(name, idNumber) {
            saveToAuthTable(addr, "idcard", newIdCard(name, idNumber))
            rest.PostProcessResponse(w, cliCtx, "Success")
        } else {
            rest.WriteErrorResponse(w, http.StatusNotFound, "Failed to verify")
        }
    }
}

func oracleVerifyCorpInfo(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        accessCode := vars["accessToken"]
        corpName   := vars["corp_name"]
        regNumber  := vars["reg_number"]
        creditCode := vars["credit_code"]

        addr, err := utils.VerifyAccessCode(accessCode)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }

        idNumber, err := getIdNumber(addr)
        if err != nil {
            rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
            return
        }

        if verifyCorpInfo(idNumber, corpName, regNumber, creditCode) {
            saveToAuthTable(addr, "corp", newCorpInfo(corpName, regNumber, creditCode))
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
    msgs := oracle.GetInsertRowMsgs("0000000001", "authentication", []types.RowFields{rowFields})

    accNum, _, err := oracle.GetAccountInfo(addr.String())
    if err == nil && accNum == 0 {
        //oracle.SendFirstTokenTo(addr)
        msg, err := oracle.GetSendTokenMsg(addr)
        if err == nil {
            msgs = append(msgs, msg)
        }
    }

    oracle.BuildTxsAndBroadcast(msgs)
}

func verifyCorpInfo(idNumber, corpName, regNumber, creditCode string) bool {
    bodyBytes, err := getPersonalBusinessData(idNumber)
    if err != nil {
        return false
    }

    regN, creditC, err := getCorpInfo(bodyBytes, corpName)
    if err != nil {
        return false
    }

    if regN == regNumber && creditC == creditCode {
        return true
    }
    return false
}

func getCorpInfo(bodyBytes []byte, corpName string) (string, string, error) {
    var ent interface{}
    err := json.Unmarshal(bodyBytes, &ent)
    if err != nil {
        return "", "", errors.New("Failed to unmarshal body bytes")
    }
    var queue = []string{"data", "data", "ryposfrs"}
    for _, v := range queue {
        ent = oracle.UnWrap(ent, v)
    }

    corps := ent.([]interface{})
    if len(corps) < 1 {
        return "", "", errors.New("No corps found")
    }

    for _, corp := range corps {
        c := corp.(map[string]interface{})
        name, ok := c["entname"]
        if !ok {
            return "", "", errors.New("Couldn't find ryname")
        }
        if name == corpName {
            if r, ok := c["regno"]; ok {
                if c, ok := c["creditcode"]; ok {
                    return  r.(string), c.(string), nil
                }
            }
        }
    }
    return "", "", errors.New("Not found!!!")
}

func verifyNameAndIdNumber(name, id_number string) bool {
    bodyBytes, err := getPersonalBusinessData(id_number)
    if err != nil {
        return false
    }
    frName, err := getName(bodyBytes)    
    if err != nil {
        fmt.Printf("Failed to get name")
        return false
    }
    if name != frName {
        fmt.Printf("Names don't match")
        return false
    }
    return true
}

func getPersonalBusinessData(idNumber string) ([]byte, error) {
    url := fmt.Sprintf("https://personalbusiness.shumaidata.com/getPersonalBusiness?idcard=%s", idNumber)
    request, err := http.NewRequest("GET", url, nil)
    if err != nil {
        //TODO: convert the println to log
        fmt.Println("Failed to create new request!!!")
        return nil, errors.New("Failed to create new request!!!")
    }

    appCode := viper.GetString(AliyunSkyEyeAppCode)
    request.Header.Set("Authorization", "APPCODE " + appCode)

    var client = &http.Client{
        Timeout: time.Second * 10,
    }
    response, err := client.Do(request)
    if err != nil {
        fmt.Println("Failed to do request!!!")
        return nil, errors.New("Failed to do request!!!")
    }

    defer response.Body.Close()
    if response.StatusCode != 200 {
        fmt.Printf("Returned code is %d!!!\n", response.StatusCode)
        return nil, errors.New(fmt.Sprintf("Returned code is %d!!!\n", response.StatusCode))
    }

    bodyBytes, _ := ioutil.ReadAll(response.Body)
    return bodyBytes, nil
}

func getName(bodyBytes []byte) (string, error) {
    var ent interface{}
    err := json.Unmarshal(bodyBytes, &ent)
    if err != nil {
        return "", errors.New("Failed to unmarshal body bytes")
    }
    var queue = []string{"data", "data", "ryposfrs"}
    for _, v := range queue {
        ent = oracle.UnWrap(ent, v)
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

func getIdNumber(addr sdk.AccAddress) (string, error) {
    querierClient := oracle.NewSuperQuerierClient("0000000001")
    querierClient.Table("authentication")
    querierClient.Equal("address", addr.String())
    querierClient.Equal("type", "idcard")
    querierClient.Last()

    records, err := querierClient.Execute()
    if err != nil {
        return "", err
    }
    jsonValue := oracle.UnWrap(records[0], "value")
    result := IdCard{}
    err = json.Unmarshal([]byte(jsonValue.(string)), &result)
    if err != nil {
        return "", errors.New("Failed to unmarshal IdCard")
    }
    return result.IdNumber, nil
}
