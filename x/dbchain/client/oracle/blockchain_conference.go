package oracle

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/afocus/captcha"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"github.com/mr-tron/base58"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/yzhanginwa/dbchain/x/dbchain/client/oracle/oracle"
	"image/color"
	"image/png"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)
// tableName conference_personal_register
// fields name mobile addr
// tableName conference_company_register
// fields corporate_name responsible position mobile addr

var conferenceAppCode string
const (
	conferenceRegister = "conference_personal_register"
	conferenceCompanyRegister = "conference_company_register"
)
var cap *captcha.Captcha

func init() {
	OracleHome := "$HOME/.dbchainoracle"
	DefaultOracleHome := os.ExpandEnv(OracleHome)
	cap = captcha.New()
	if err := cap.SetFont(DefaultOracleHome + "/config/comic.ttf"); err != nil {
		panic(err.Error())
	}
	cap.SetSize(128, 64)
	cap.SetDisturbance(captcha.NORMAL)
	cap.SetFrontColor(color.RGBA{255, 255, 255, 255})
	cap.SetBkgColor(color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 153, 0, 255})
}

func oracleSendPictureVerifyCode(cliCtx context.CLIContext) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {

		img, str := cap.Create(4, captcha.UPPER)
		buf := make([]byte,0 )
		write := bytes.NewBuffer(buf)
		png.Encode(write, img)
		baseStr := base64.StdEncoding.EncodeToString(write.Bytes())
		cacheMobileAndVerificationCode(str, str, str)
		result := map[string]string {
			"picture" : baseStr,
		}
		re,_ := json.Marshal(result)
		successResponse(w,re)
	}
}

func oracleConferencePersonalRegister(cliCtx context.CLIContext) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		appType := vars["app_type"] // appType Only can applet or h5
		wechatId := vars["wechat_id"]
		name     := vars["name"]
		mobile     := vars["mobile"]
		verifyCode     := strings.ToUpper(vars["verify_code"])

		fieldValue := map[string]string {
			"wechat_id" : wechatId,
			"name" : name,
			"mobile" : mobile,
		}
		registerCore(cliCtx, w, appType, verifyCode, conferenceRegister, fieldValue)
	}
}

func oracleConferenceCorporateRegister(cliCtx context.CLIContext) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		wechatId := vars["wechat_id"]
		corporateName := vars["company_name"]
		responsible := vars["responsible"]
		position := vars["position"]
		mobile := vars["mobile"]
		appType := vars["app_type"] // appType Only can applet or h5
		inviter := vars["inviter"]
		verifyCode     := strings.ToUpper(vars["verify_code"])

		fieldValue := map[string]string {
			"wechat_id" : wechatId,
			"corporate_name" : corporateName,
			"responsible" : responsible,
			"position" : position,
			"inviter" : inviter,
			"mobile" : mobile,
		}
		registerCore(cliCtx, w, appType, verifyCode, conferenceCompanyRegister, fieldValue)
	}
}


func showConferenceRegistrationStatus(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		appCode := loadConferenceAppCode()
		storeName := "dbchain"
		tableName := []string{conferenceCompanyRegister, conferenceRegister}

		params := vars["params"]
		var result []map[string]string
		var err error
		for _, table := range tableName {
			result, err = getConferenceRegistrationStatus(cliCtx, storeName, appCode, table, params)
			if len(result) != 0 {
				break
			}
		}

		if err != nil || len(result) == 0 {
			generalResponse(w, map[string]string {
				"error" : "unregistered",
			})
			return
		}

		data := make(map[string]string, 0)
		for key, val := range result[0] {
			if key == "id" || key == "created_by" || key == "created_at" {
				continue
			}
			data[key] = val
		}
		bz,_ := json.Marshal(data)
		successResponse(w,bz)
	}
}

func getConferenceRegistrationStatus(cliCtx context.CLIContext, storeName, appCode, tableName, params string) ([]map[string]string, error) {
	fields := make(map[string]string)
	fields["wechat_id"] = params

	result , err := queryByWhere(cliCtx, storeName, appCode, tableName, fields)
	return result, err
}
func showConferenceRegisterNumbers(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ac := getOracleAc()
		appCode := loadConferenceAppCode()
		tableName := conferenceRegister
		userType := vars["user_type"]
		if userType == "company_name" {
			tableName = conferenceCompanyRegister
		}
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/find_all/%s/%s/%s", storeName, ac, appCode, tableName), nil)
		if err != nil {
			generalResponse(w, map[string]string {
				"error" : err.Error(),
			})
			return
		}
		ids := make([]string, 0)
		err = json.Unmarshal(res,&ids)
		if err != nil {
			generalResponse(w, map[string]string {
				"error" : err.Error(),
			})
			return
		}
		data := map[string]string{
			"registers" : fmt.Sprintf("%d",len(ids)),
		}
		bz,_ := json.Marshal(data)
		successResponse(w,bz)
	}
}

func loadConferenceAppCode() string{
	if conferenceAppCode != ""{
		return conferenceAppCode
	}
	conferenceAppCode = viper.GetString("conference-appcode")
	return conferenceAppCode
}

func generalResponse(w http.ResponseWriter, data interface{}) {
	bz,_ := json.Marshal(data)
	successResponse(w,bz)
}

func getWeChatUserInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		result, err  := ioutil.ReadAll(r.Body)
		if err != nil {
			generalResponse(w, map[string]string{"error" : err.Error()})
			return
		}
		data := make(map[string]string)
		err = json.Unmarshal(result, &data)
		if err != nil {
			generalResponse(w, map[string]string{"error" : err.Error()})
			return
		}

		postUrl := "https://api.weixin.qq.com/sns/jscode2session"
		DataUrlVal := url.Values{}
		for key,val := range data{
			DataUrlVal.Add(key,val)
		}
		contentType := "application/x-www-form-urlencoded"
		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Post(postUrl, contentType, bytes.NewBuffer([]byte(DataUrlVal.Encode())))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		result, _ = ioutil.ReadAll(resp.Body)
		successResponse(w,result)
		return
	}
}

func registerCore(cliCtx context.CLIContext, w http.ResponseWriter, appType, verifyCode string, tableName string, fieldValue map[string]string) {
	if appType != "h5" && appType != "applet" {
		generalResponse(w, map[string]string {
			"error" : "invalid request",
		})
		return
	}
	if appType == "h5" && !VerifyVerfCode(verifyCode, verifyCode, verifyCode) {
		generalResponse(w, map[string]string {
			"error" : "verify code error",
		})
		return
	}
	mobile := fieldValue["mobile"]
	//if register is corporate, add corporateName on mobile
	if corporateName, ok := fieldValue["corporate_name"]; ok {
		mobile = corporateName + "_" + mobile
	}
	//Check for duplicate registration
	if IsCachedMobileAndVerificationCode(mobile) {
		generalResponse(w, map[string]string {
			"error" : "repeat registration",
		})
		return
	}

	pk := secp256k1.GenPrivKey()
	addr := sdk.AccAddress(pk.PubKey().Address())
	fieldValue["addr"] = addr.String()
	//can insert

	jsFieldValue,_  := json.Marshal(fieldValue)
	base58FieldValue := base58.Encode(jsFieldValue)
	ac := getOracleAc()
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/can_insert_row/%s/%s/%s/%s","dbchain", ac, loadConferenceAppCode(), tableName, base58FieldValue), nil)
	if err != nil {
		generalResponse(w, map[string]string {
			"error" : err.Error(),
		})
		return
	}
	var canInsert bool
	err = json.Unmarshal(res, &canInsert)
	if err != nil {
		generalResponse(w, map[string]string {
			"error" : err.Error(),
		})
		return
	}
	if !canInsert {
		generalResponse(w, map[string]string {
			"error" : "this telephone number has been registered",
		})
		return
	}
	//
	data := map[string]string{
		"register_code" : addr.String(),
		"register_status" : "success",
	}
	bz,_ := json.Marshal(data)
	successResponse(w,bz)
	// Cache registered phones
	cacheMobileAndVerificationCode(mobile,mobile,"")
	oracleAccAddr := oracle.GetOracleAccAddr()
	OracleSaveToTable(cliCtx, oracleAccAddr, fieldValue, tableName, loadConferenceAppCode())
}