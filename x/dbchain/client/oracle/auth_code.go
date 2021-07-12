package oracle

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/yzhanginwa/dbchain/x/dbchain/client/oracle/authenticator"
	"github.com/yzhanginwa/dbchain/x/dbchain/client/oracle/oracle"
	"github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	 authCodeInfo = "auth_code_info"
	 //
	organizationUserSecretKey = "organization_user_secret_key"
	secretKeyAuth = "secret_key_authentication"
)

func organizationGetSecretKey(cliCtx context.CLIContext) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()
		organizationName := r.Form.Get("organization_name")
		organizationId := r.Form.Get("organization_id")
		//TODO  query form database
		fieldValue := map[string]string {
			"organization_name" :  organizationName,
			"organization_id" : organizationId,
			"type" : "organization",
		}
		//1. query from database
		secrets, err := queryByWhere(cliCtx, "dbchain", "0000000001", organizationUserSecretKey, fieldValue)
		if len(secrets) != 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "organization has registered")
			return
		}
		//2. generate key
		secret , encryptSecret, err := generateSecretKey(cliCtx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, "generate secret key err")
			return
		}
		//save to table
		fieldValue["secret_key"] = encryptSecret
		oracleAccAddr := oracle.GetOracleAccAddr()
		SaveToOrderInfoTable(cliCtx, oracleAccAddr, fieldValue, organizationUserSecretKey)
		//TODO need conform write success before return
		b := queryTxResult(cliCtx, "dbchain", "0000000001", organizationUserSecretKey, fieldValue)
		if b {
			result := map[string]string {
				"secret_key" : secret,
			}
			rest.PostProcessResponse(w, cliCtx, result)
		} else {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, "generate secret key err")
		}
	}
}

//
func userGetSecretKey(cliCtx context.CLIContext) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		organizationId := r.Form.Get("organization_id")
		//query organization_name
		queryOrganizationName := map[string]string {
			"organization_id" : organizationId,
			"type" : "organization",
		}
		//只有通过验证的平台的用户才是合法用户
		or, err := queryByWhere(cliCtx, "dbchain", "0000000001", organizationUserSecretKey, queryOrganizationName)
		if err != nil || len(or) == 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid user")
			return
		}
		orId := or[len(or)-1]["id"]
		if !isAuthentication(cliCtx, orId) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid organization")
			return
		}
		//
		organizationName := or[len(or)-1]["organization_name"]
		userId := r.Form.Get("user_id")
		//
		fieldValue := map[string]string {
			"user_id" :  userId,
			"organization_name" : organizationName,
			"organization_id" : organizationId,
			"type" : "user",
		}
		//1. if it has register  return err
		secrets, err := queryByWhere(cliCtx, "dbchain", "0000000001", organizationUserSecretKey, fieldValue)
		if len(secrets) != 0 {
			ID := secrets[len(secrets)-1]["id"]
			queryAuth := map[string]string {
				"organization_user_secret_key_id" : ID,
			}
			auth, err := queryByWhere(cliCtx, "dbchain", "0000000001", secretKeyAuth, queryAuth)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, "service err")
				return
			}
			if  len(auth) != 0 {
				rest.WriteErrorResponse(w, http.StatusBadRequest, "user has authed")
				return
			}
		}

		//generate secret key
		secret ,encryptSecret,  err := generateSecretKey(cliCtx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, "generate secret key err")
			return
		}
		fieldValue["secret_key"] = encryptSecret
		oracleAccAddr := oracle.GetOracleAccAddr()
		SaveToOrderInfoTable(cliCtx, oracleAccAddr, fieldValue, organizationUserSecretKey)

		b := queryTxResult(cliCtx, "dbchain", "0000000001", organizationUserSecretKey, fieldValue)
		if b {
			data , err := genQrCodeString(secret, organizationName, userId)
			if err !=nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, "generate secret key err")
				return
			}
			//rest.PostProcessResponse(w, cliCtx, data)
			successResponse(w, []byte(data))
		} else {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, "generate secret key err")
			return
		}
	}
}

func organizationVerify(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		state := r.Form.Get("state")
		organizationId := r.Form.Get("organization_id")
		hashKey := r.Form.Get("hash_key")
		queryOrganizationName := map[string]string {
			"organization_id" : organizationId,
			"type" : "organization",
		}
		or, err := queryByWhere(cliCtx, "dbchain", "0000000001", organizationUserSecretKey, queryOrganizationName)
		if err != nil || len(or) == 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid organization")
			return
		}
		encryptSecretKey := or[len(or)-1]["secret_key"]
		ID := or[len(or)-1]["id"]
		secretKey, err := decryptSecretKey(encryptSecretKey)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, "service err")
			return
		}
		hs := sha256.Sum256([]byte(secretKey + state))
		haHex := hex.EncodeToString(hs[:])
		if strings.ToUpper(haHex) != strings.ToUpper(hashKey) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid hash_key")
		} else {
			result := "success"
			reHs := sha256.Sum256([]byte(secretKey + state + result))
			reHsHex := hex.EncodeToString(reHs[:])
			data := map[string]string {
				"hash_key" : reHsHex,
				"result" : result,
				"state" : state,
			}
			//
			fieldValue := map[string]string {
				"organization_user_secret_key_id" : ID,
				"auth_status" : "true",
			}
			oracleAccAddr := oracle.GetOracleAccAddr()
			SaveToOrderInfoTable(cliCtx, oracleAccAddr, fieldValue, secretKeyAuth)
			b := queryTxResult(cliCtx, "dbchain", "0000000001", secretKeyAuth, fieldValue)
			if !b {
				data["result"] = "field"
			}
			bz , _ := json.Marshal(data)
			successResponse(w, bz)
		}

	}
}


func userVerifyCode(cliCtx context.CLIContext, enableAuth bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		state := r.Form.Get("state")
		organizationId := r.Form.Get("organization_id")
		userId := r.Form.Get("user_id")
		verifyCode := r.Form.Get("verify_code")

		//
		queryOrganizationName := map[string]string {
			"user_id" : userId,
			"type" : "user",
			"organization_id" : organizationId,
		}
		or, err := queryByWhere(cliCtx, "dbchain", "0000000001", organizationUserSecretKey, queryOrganizationName)
		if err != nil || len(or) == 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid user")
			return
		}
		encryptSecretKey := or[len(or)-1]["secret_key"]
		ID := or[len(or)-1]["id"]
		secretKey, err := decryptSecretKey(encryptSecretKey)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, "service err")
			return
		}
		//need user has authed
		if enableAuth == false {
			queryAuth := map[string]string {
				"organization_user_secret_key_id" : ID,
			}
			or, err := queryByWhere(cliCtx, "dbchain", "0000000001", secretKeyAuth, queryAuth)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, "service err")
				return
			}
			if  len(or) == 0 {
				rest.WriteErrorResponse(w, http.StatusBadRequest, "user dont auth")
				return
			}
		}

		ga := authenticator.NewGAuth()
		ret, err := ga.VerifyCode(secretKey, verifyCode, 1)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, "service err")
		}

		//
		result := "success"
		if ret == false {
			result = "field"
		}
		reHs := sha256.Sum256([]byte(secretKey + state + result))
		reHsHex := hex.EncodeToString(reHs[:])
		//
		if enableAuth && result == "success" {
			fieldValue := map[string]string {
				"organization_user_secret_key_id" : ID,
				"auth_status" : "true",
			}
			oracleAccAddr := oracle.GetOracleAccAddr()
			SaveToOrderInfoTable(cliCtx, oracleAccAddr, fieldValue, secretKeyAuth)
			b := queryTxResult(cliCtx, "dbchain", "0000000001", secretKeyAuth, fieldValue)
			if !b {
				result = "field"
			}
		}
		//redirect
		data := map[string]string{
			"state" : state,
			"result" : result,
			"user_id" : userId,
			"hash_key" : reHsHex,
		}
		bz,_ := json.Marshal(data)
		successResponse(w,bz)
		return
	}
}

func userDestoryCode(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		organizationId := r.Form.Get("organization_id")
		userId := r.Form.Get("user_id")
		state := r.Form.Get("state")

		queryOrganizationName := map[string]string {
			"user_id" : userId,
			"type" : "user",
			"organization_id" : organizationId,
		}

		or, err := queryByWhere(cliCtx, "dbchain", "0000000001", organizationUserSecretKey, queryOrganizationName)
		if err != nil || len(or) == 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid user")
			return
		}

		encryptSecretKey := or[len(or)-1]["secret_key"]
		ID := or[len(or)-1]["id"]
		secretKey, err := decryptSecretKey(encryptSecretKey)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, "service err")
			return
		}

		queryAuth := map[string]string {
			"organization_user_secret_key_id" : ID,
		}
		auth, err := queryByWhere(cliCtx, "dbchain", "0000000001", secretKeyAuth, queryAuth)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, "service err")
			return
		}
		if  len(auth) == 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "user has not authed")
			return
		}
		oracleAccAddr := oracle.GetOracleAccAddr()
		msgs := make([]oracle.UniversalMsg, 0)
		//
		for i := 0; i < len(or); i++ {
			sId := or[i]["id"]
			id , _ := strconv.Atoi(sId)
			msg := types.NewMsgFreezeRow(oracleAccAddr, "0000000001", organizationUserSecretKey, uint(id))
			msgs = append(msgs, msg)
			//
			queryAuth := map[string]string {
				"organization_user_secret_key_id" : sId,
			}
			auth, err := queryByWhere(cliCtx, "dbchain", "0000000001", secretKeyAuth, queryAuth)
			if err != nil || len( auth) == 0 {
				continue
			}
			 for j := 0; j < len(auth); j++ {
				 sId := auth[j]["id"]
				 id , _ := strconv.Atoi(sId)
				 msg := types.NewMsgFreezeRow(oracleAccAddr, "0000000001", secretKeyAuth, uint(id))
				 msgs = append(msgs, msg)
			 }
		}

		oracle.BuildTxsAndBroadcast(cliCtx, msgs)
		result := "success"
		b := queryFreezeResult(cliCtx, "dbchain", "0000000001", secretKeyAuth, queryAuth)
		if !b {
			result = "field"
		}
		//redirect

		reHs := sha256.Sum256([]byte(secretKey + state + result))
		reHsHex := hex.EncodeToString(reHs[:])
		data := map[string]string{
			"state" : state,
			"result" : result,
			"user_id" : userId,
			"hash_key" : reHsHex,
		}
		bz,_ := json.Marshal(data)
		successResponse(w,bz)
		return
	}
}

func queryTxResult(cliCtx context.CLIContext, storeName, appcode, tableName string , fieldValue map[string]string) bool {
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		secrets, _ := queryByWhere(cliCtx, storeName, appcode, tableName, fieldValue)
		if len(secrets) > 0 {
			return true
		}
	}
	return false
}

func generateSecretKey(cliCtx context.CLIContext) (string,string, error) {
	for i := 0; i < 10; i++ {
		ga := authenticator.NewGAuth()
		secret , err := ga.CreateSecret(32)
		if err != nil {
			continue
		}
		encryptSecret, err := encryptedSecretKey(secret)
		if err != nil {
			continue
		}
		fieldValue := map[string]string {
			"secret_key" : encryptSecret,
		}
		secrets, err := queryByWhere(cliCtx, "dbchain", "0000000001", organizationUserSecretKey, fieldValue)
		if len(secrets) != 0 {
			continue
		}
		return secret, encryptSecret, nil
	}

	return "", "", errors.New("generate secret key")
}

func encryptedSecretKey(secretKey string) (string,error) {
	encryptKey , err := getAesKey()
	if err != nil {
		return "", err
	}
	encryptData,err := AESEncrypt([]byte(secretKey),encryptKey)
	if err != nil {
		return "", nil
	}
	encryptString := hex.EncodeToString(encryptData)
	return encryptString, nil
}

func decryptSecretKey(encryptedSecretKey string) (string, error) {
	encryptKey , err := getAesKey()
	if err != nil {
		return "", err
	}
	encryptedSecretKeyData , _ := hex.DecodeString(encryptedSecretKey)
	secretKeyData, err := AESDecrypt(encryptedSecretKeyData, encryptKey)
	if err != nil {
		return "", nil
	}
	return string(secretKeyData), nil
}

// AES加密
 func AESEncrypt(src, key []byte) ([]byte,error){
	 block, err := aes.NewCipher(key)
	 if err != nil{
	 	return nil, err
	 }
	 src = PKCS5Padding(src, block.BlockSize())
	 blockMode := cipher.NewCBCEncrypter(block, key[:block.BlockSize()])
	 dst := src
	 blockMode.CryptBlocks(dst, src)
	return dst,nil

}

// AES解密
 func AESDecrypt(src, key []byte) ([]byte,error){
	 block, err := aes.NewCipher(key)
	 if err != nil{
		return nil, err
	 }
	 blockMode := cipher.NewCBCDecrypter(block, key[:block.BlockSize()])
	 dst := src
	 blockMode.CryptBlocks(dst, src)
	 dst = PKCS5UnPadding(dst)
	 return dst,nil
}

// 使用pks5的方式填充
func PKCS5Padding(ciphertext []byte, blockSize int) []byte{
	padding := blockSize - (len(ciphertext)%blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	newText := append(ciphertext, padText...)
	return newText
}

// 删除pks5填充的尾部数据
func PKCS5UnPadding(origData []byte) []byte{
	length := len(origData)
	number := int(origData[length-1])
	return origData[:(length-number)]
}

func getAesKey() ([]byte,error) {
	//采用 aes 加密，密钥为 Oracle privateKey hash
	privKey, err := oracle.LoadPrivKey()
	if err != nil {
		return nil, err
	}
	hashData := sha256.Sum256(privKey[:])
	return hashData[:] ,nil
}

func genQrCodeString(secretKey, organizationName ,userId string) (string,error) {
	//"otpauth://totp/Firefox:893006326%40qq.com?secret=OJYVS3LFNBBTQ4TNOBFGCR2YKJYFMV2R&issuer=Firefox"
	userId = url.QueryEscape(userId)
	QrCodeString := fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s", organizationName, userId, secretKey, organizationName)
	q, err := qrcode.New(QrCodeString, qrcode.Medium)
	if err != nil {
		return "", err
	}
	bStream, err  := q.PNG(200)
	if err != nil {
		return "", nil
	}

	QrBase := base64.StdEncoding.EncodeToString(bStream)


	obj := map[string]string{
		"qrcode_url_base64" : QrBase,
		"secret_key" : secretKey,
	}
	js ,err := json.Marshal(obj)
	if err != nil {
		return "", nil
	}
	return string(js), nil
}

func successResponse(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}

func queryFreezeResult(cliCtx context.CLIContext, storeName, appcode, tableName string , fieldValue map[string]string ) bool {
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		auth, err := queryByWhere(cliCtx, storeName, appcode, tableName, fieldValue)
		if err != nil {
			return false
		}
		if  len(auth) == 0 {
			return true
		}
	}
	return false

}

func isAuthentication( cliCtx context.CLIContext, ID string) bool {
	queryAuth := map[string]string {
		"organization_user_secret_key_id" : ID,
	}
	auth, err := queryByWhere(cliCtx, "dbchain", "0000000001", secretKeyAuth, queryAuth)
	if err != nil || len(auth) == 0{
		return false
	}
	return true
}
