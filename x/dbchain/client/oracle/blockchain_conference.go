package oracle

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
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
	"math/big"
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
	euCodeOfPersonalRegister = "eucode_of_personal_register"
)
var cap *captcha.Captcha

func loadCap() error {
	OracleHome := "$HOME/.dbchainoracle"
	DefaultOracleHome := os.ExpandEnv(OracleHome)
	cap = captcha.New()
	if err := cap.SetFont(DefaultOracleHome + "/config/comic.ttf"); err != nil {
		fmt.Println(err.Error())
		return err
	}
	cap.SetSize(128, 64)
	cap.SetDisturbance(captcha.NORMAL)
	cap.SetFrontColor(color.RGBA{255, 255, 255, 255})
	cap.SetBkgColor(color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 153, 0, 255})
	return nil
}

func oracleSendPictureVerifyCode(cliCtx context.CLIContext) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {

		if cap == nil {
			err := loadCap()
			if err != nil {
				successResponse(w, nil)
				return
			}
		}
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
		uiCardidType := vars["ui_cardid_type"]
		uiCardid := vars["ui_cardid"]

		if checkIfRegister(cliCtx, wechatId) {
			generalResponse(w, map[string]string {
				"error" : "this telephone number has been registered",
			})
			return
		}

		euCode, err := oracleConferenceAuthentication(cliCtx, uiCardidType, uiCardid, name)
		if err != nil {
			generalResponse(w, map[string]string {
				"error" : err.Error(),
			})
			return
		}
		err = oracleConferenceUpdateRegisterIdentityInfo(cliCtx, wechatId, euCode)
		if err != nil {
			generalResponse(w, map[string]string {
				"error" : err.Error(),
			})
			return
		}
		fieldValue := map[string]string {
			"wechat_id" : wechatId,
			"name" : name,
			"mobile" : mobile,
		}
		registerCore(cliCtx, w, appType, verifyCode, conferenceRegister, fieldValue)
	}
}

func checkIfRegister(cliCtx context.CLIContext, wechatId string) bool {
	ac := getOracleAc()
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/find_by/%s/%s/%s/%s/%s","dbchain", ac, loadConferenceAppCode(), conferenceRegister, "wechat_id", wechatId), nil)
	if err != nil {
		return true
	}
	ids := make([]string, 0)
	json.Unmarshal(res, &ids)
	if len(ids) == 0 {
		return false
	}
	return true
}
func oracleConferenceAuthentication(cliCtx context.CLIContext, uiCardidType, uiCardid, name string)  (string,error) {
	code := genIdentityCode(cliCtx)
	if code == "" {
		return "", errors.New("gen euCode code err")
	}

	data := map[string]string {
		"appKey": "38DB53D9742B7A1D398EBC606FF4C1365127EF8A185E108220428F0219EBF6E3",
		"appScret": "38DB53D9742B7A1D398EBC606FF4C1365127EF8A185E1082FC627FE6216C6E6B9D7D7FA4EDA09946",
		"euCode": code,
		"uiName": name,
		"uiCardidType" : uiCardidType,
		"uiCardid" : uiCardid,
	}

	bz, _ := json.Marshal(data)
	buf := bytes.NewReader(bz)
	resp, err := http.Post("http://reg.dataexpo.com.cn/gate-api/api/accept.do", "application/json", buf)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	result := make(map[string]interface{})
	err = json.Unmarshal(respData, &result)
	if err != nil {
		return "", err
	}
	ErrCode := int(result["errcode"].(float64))
	ErrMsg := result["errmsg"].(string)


	if (ErrCode == 0 && ErrMsg == "成功") || (ErrCode == 1 && ErrMsg == "证件码已存在") {
		return code,nil
	}

	return "", errors.New(ErrMsg)
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

func genIdentityCode(cliCtx context.CLIContext) string {
	max := big.NewInt(89999999)
	high := big.NewInt(10000000)
	for i := 0; i < 20; i++ {
		euCode, err := rand.Int(rand.Reader, max)
		if err != nil {
			continue
		}
		euCode.Add(euCode, high)
		ac := getOracleAc()
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/find_by/%s/%s/%s/%s/%s","dbchain", ac, loadConferenceAppCode(), euCodeOfPersonalRegister, "eu_code", euCode.String()), nil)
		if err != nil {
			continue
		}
		ids := make([]string, 0)
		json.Unmarshal(res, &ids)
		if len(ids) == 0 {
			return euCode.String()
		}
	}
	return ""
}

func oracleConferenceUpdateRegisterIdentityInfo(cliCtx context.CLIContext, wechatId, euCode string) error{

	fieldValue := map[string]string {
		"wechat_id" : wechatId,
		"eu_code" : euCode,

	}
	oracleAccAddr := oracle.GetOracleAccAddr()
	OracleSaveToTable(cliCtx, oracleAccAddr, fieldValue, euCodeOfPersonalRegister, loadConferenceAppCode())
	return nil
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
		identity := getRegisterIdentity(cliCtx, storeName, appCode, euCodeOfPersonalRegister, params)
		if identity != "" {
			data["identity"] =  identity
		} else {
			data["identity"] =  fmt.Sprintf("%08s", result[0]["id"])
		}

		bz,_ := json.Marshal(data)
		successResponse(w,bz)
	}
}

func getRegisterIdentity(cliCtx context.CLIContext, storeName, appCode, table, wechatId string) string{
	result, _ := getConferenceRegistrationStatus(cliCtx, storeName, appCode, table, wechatId)
	if len(result) != 0 {
		return result[0]["eu_code"]
	}
	return ""
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