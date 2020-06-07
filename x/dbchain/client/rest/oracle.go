package rest

import (
    "fmt"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/spf13/viper"
    "github.com/mr-tron/base58"
    "github.com/tendermint/tendermint/crypto/secp256k1"
    "github.com/cosmos/cosmos-sdk/client/context"
    "github.com/cosmos/cosmos-sdk/types/rest"
    sdk "github.com/cosmos/cosmos-sdk/types"
    authtypes "github.com/cosmos/cosmos-sdk/x/auth/types
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
)

type MobileVerfCode struct {
    Mobile string     `json:"mobile"`
    VerfCode string   `json:"verf_code"`
}

const OracleEncryptedPrivKey = "oracle-encrypted-key"

var (
    associationMap = make(map[string]MobileVerfCode)
    oraclePrivKey secp256k1.PrivKeySecp256k1
    oraclePrivKeyLoaded = false
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
    fmt.Println(mobile)
    fmt.Println(verificationCode)
    a := viper.GetString("ethan")
    fmt.Println(a)
    viper.Set("ethan", "a new value")
    viper.SafeWriteConfig()
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

    stdSignMsg := authtypes.StdSignMsg{
        ChainID:       "mainchain",
        AccountNumber: 1,
        Sequence:      1,
        Memo:          "hello from Voyager 1!",
        Msgs:          msgs,
        Fee:           authtypes.NewStdFee(200000, sdk.Coins{sdk.NewCoin("dbctoken", sdk.NewInt(1))}),
    },

    sig, err = priv.Sign(stdSignMsg.Bytes())
    if err != nil {
        return nil, nil, err
    }

    return sig, priv.PubKey(), nil

    signBytes = sig

    return authtypes.StdSignature{
        PubKey:    pubkey,
        Signature: sigBytes,
    }


	func NewStdTx(msgs []sdk.Msg, fee StdFee, sigs []StdSignature, memo string) StdTx {
    33:		return StdTx{
=>  34:			Msgs:       msgs,
    35:			Fee:        fee,
    36:			Signatures: sigs,
    37:			Memo:       memo,
    38:		}
    39:	}









}
