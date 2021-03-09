package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"math/big"
)


const (
	base = "wzAG2dsfbkrEinDKPamBpQ6WtUuHLNceyRVXZ78h3TCJSY5qxjvM14"
	decimal = 54
	pad = "F"
	codeLen = 6
)

func numTOCode(num int64) string {
	mod := int64(0)
	res := ""

	for num!=0 {
		mod = num % decimal
		num = num / decimal
		res += string(base[mod])
	}
	resLen := len(res)
	if resLen < codeLen {
		res += pad
		randSeed := big.NewInt(decimal)
		for i:=0; i< codeLen - resLen - 1; i++ {
			Bi, _ := rand.Int(rand.Reader,randSeed)
			res += string(base[Bi.Int64()])
		}
	}
	return res
}

func GenInvitationCod(cliCtx context.CLIContext, appCode, tableName, addr string) string {
	originData := []byte(addr)
	for {
		hashData := sha256.Sum256(originData)
		reader := bytes.NewReader(hashData[:])
		num , err := rand.Prime(reader, 32)
		if err != nil {
			originData = hashData[:]
			continue
		}
		invitationCode := numTOCode(num.Int64())
		if checkCode(cliCtx, appCode, tableName, invitationCode) {
			return invitationCode
		} else {
			originData = hashData[:]
		}
	}
}

func checkCode(cliCtx context.CLIContext, appCode, tableName, invitationCode string) bool {
	privKey, err := LoadPrivKey()
	if err != nil {
		panic("failed to load oracle's private key!!!")
	}

	accessToken := MakeAccessCode(privKey)
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/find_by/%s/%s/%s/%s/%s", "dbchain", accessToken, appCode, tableName, "code", invitationCode), nil)
	if err != nil {
		return false
	}
	var out []string
	json.Unmarshal(res,&out)
	if len(out) == 0 {
		return true
	}

	return false
}