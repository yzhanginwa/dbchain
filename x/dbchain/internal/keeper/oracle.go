package keeper

import (
    "sort"
    "encoding/json"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/client/rest/oracle"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
)

type OracleAuthFields struct {
    Type string
    Mobile string
    Name string
    IdNumber string
    CorpName string
    RegNumber string
    CreditCode string
}

func checkWithOracleAuth(keeper Keeper, ctx sdk.Context, fields types.RowFields, owner sdk.AccAddress) bool {
    authType, ok := fields["type"]
    if !ok { return false }

    records, err := getOracleAuthRecords(keeper, ctx, owner, authType)
    if err != nil {
        return false
    }

    if len(records) < 1 {
        return false
    }

    oracleAuthFields := importAuthValueIntoAuthFields(records)

    if len(oracleAuthFields) < 1 {
        return false
    }
    authRecord := oracleAuthFields[0]

    switch authType {
    case "mobile":
        if fields["mobile"] == authRecord.Mobile {
            return true
        }
        return false
    case "idcard":
        if fields["name"] == authRecord.Name &&
           fields["id_number"] == authRecord.IdNumber {
            return true
        }
        return false
    case "corp":
        if fields["corp_name"]   == authRecord.CorpName &&
           fields["reg_number"]  == authRecord.RegNumber &&
           fields["credit_code"] == authRecord.CreditCode {
            return true
        }
        return false
    default:
        return false
    }
    return true
}


func getOracleAuthRecords(keeper Keeper, ctx sdk.Context, owner sdk.AccAddress, authType string) ([]map[string]string, error) {
    appId, err := keeper.GetDatabaseId(ctx, "0000000001")
    if err != nil {
        return nil, err
    }

    querierObjs := []map[string]string{}
    var ent map[string]string

    ent = map[string]string{
        "method": "table",
        "table": "authentication",
    }
    querierObjs = append(querierObjs, ent)

    ent = map[string]string{
        "method": "equal",
        "field": "address",
        "value": owner.String(),
    }
    querierObjs = append(querierObjs, ent)

    //TODO: to let the super querier support filtering out records on multiple criteria
    rows, _, err := querierSuperHandler(ctx, keeper, appId, querierObjs, owner)
    if err != nil {
        return nil, err
    }

    result := []map[string]string{}

    for _, row := range rows {
        if t, ok := row["type"]; ok {
            if t == authType {
                result = append(result, row)
            }
        }
    }

    //TODO: to let the super querier support ordering
    sort.Slice(result, func(i, j int) bool {return result[i]["id"] > result[j]["id"]})
    return result, nil
}

func importAuthValueIntoAuthFields(rows []map[string]string) []OracleAuthFields {
    result := []OracleAuthFields{}
    for _, row := range rows {
        authFields := OracleAuthFields{}
        authFields.Type = row["type"]
        value := row["value"]
        var data interface{}
        err := json.Unmarshal([]byte(value), &data)
        if err != nil {
            return nil
        }
        switch authFields.Type {
        case "mobile":
            authFields.Mobile = oracle.UnWrap(data, "mobile").(string)
        case "idcard":
            authFields.Name     = oracle.UnWrap(data, "name").(string)
            authFields.IdNumber = oracle.UnWrap(data, "id_number").(string)
        case "corp":
            authFields.CorpName   = oracle.UnWrap(data, "corp_name").(string)
            authFields.RegNumber  = oracle.UnWrap(data, "reg_number").(string)
            authFields.CreditCode = oracle.UnWrap(data, "credit_code").(string)
        }
        result = append(result, authFields)
    }
    return result
}
