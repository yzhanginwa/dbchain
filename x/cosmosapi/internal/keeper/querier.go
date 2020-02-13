package keeper

import (
    "fmt"
    "strings"
    "strconv"
    "errors"
    "time"
    "regexp"
    //"encoding/hex"
    "encoding/base64"
    "github.com/tendermint/tendermint/crypto/secp256k1"

    "github.com/cosmos/cosmos-sdk/codec"
    sdk "github.com/cosmos/cosmos-sdk/types"
    abci "github.com/tendermint/tendermint/abci/types"
    //"github.com/yzhanginwa/cosmos-api/x/cosmosapi/internal/types"
)

// query endpoints supported by the cosmosapi service Querier
const (
    QueryTables   = "tables"
    QueryIndex    = "index"
    QueryOption   = "option"
    QueryFieldOption   = "field_option"
    QueryRow      = "find"
    QueryIdsBy    = "find_by"
    QueryAdminGroup = "admin_group"

    MaxAllowedTimeDiff = 15 * 1000   // 15 seconds
)


// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
    return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
        switch path[0] {
        case QueryTables:
            if len(path) > 1 {
                return queryTable(ctx, path[1:], req, keeper)
            } else {
                return queryTables(ctx, req, keeper)
            }
        case QueryIndex:
            return queryIndex(ctx, path[1:], req, keeper)
        case QueryOption:
            return queryOption(ctx, path[1:], req, keeper)
        case QueryFieldOption:
            return queryFieldOption(ctx, path[1:], req, keeper)
        case QueryRow:
            return queryRow(ctx, path[1:], req, keeper)
        case QueryIdsBy:
            return queryIdsBy(ctx, path[1:], req, keeper)
        case QueryAdminGroup:
            return queryAdminGroup(ctx, path[1:], req, keeper)
        default:
            return nil, sdk.ErrUnknownRequest("unknown cosmosapi query endpoint")
        }
    }
}

func queryTables(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
    tables, err := keeper.getTables(ctx)

    if err != nil {
        return []byte{}, sdk.ErrUnknownRequest("Can not get table names")
    }

    res, err := codec.MarshalJSONIndent(keeper.cdc, tables)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryTable(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
    table, err := keeper.GetTable(ctx, path[0])

    if err != nil {
        return []byte{}, sdk.ErrUnknownRequest(fmt.Sprintf("Table %s does not exist",  path[0]))
    }

    res, err := codec.MarshalJSONIndent(keeper.cdc, table)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryIndex(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
    index, err := keeper.GetIndex(ctx, path[0])

    if err != nil {
        return []byte{}, sdk.ErrUnknownRequest(fmt.Sprintf("Table %s does not exist",  path[0]))
    }

    res, err := codec.MarshalJSONIndent(keeper.cdc, index)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryOption(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
    options, err := keeper.GetOption(ctx, path[0])

    if err != nil {
        return []byte{}, sdk.ErrUnknownRequest(fmt.Sprintf("Table %s does not exist",  path[0]))
    }

    res, err := codec.MarshalJSONIndent(keeper.cdc, options)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryFieldOption(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
    options, err := keeper.GetFieldOption(ctx, path[0], path[1])

    if err != nil {
        return []byte{}, sdk.ErrUnknownRequest(fmt.Sprintf("Field %s.%s does not exist",  path[0], path[1]))
    }

    res, err := codec.MarshalJSONIndent(keeper.cdc, options)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}
func queryRow(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
    u32, err := strconv.ParseUint(path[1], 10, 32)
    fields, err := keeper.Find(ctx, path[0], uint(u32))

    if err != nil {
        return []byte{}, sdk.ErrUnknownRequest(fmt.Sprintf("Table %s does not exist",  path[0]))
    }

    res, err := codec.MarshalJSONIndent(keeper.cdc, fields)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryIdsBy(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
    accessCode:= path[0]
    addr, err := verifyAccessCode(accessCode)
    fmt.Printf("Address is %s", addr)  // to make program compile

    if err != nil {
        return []byte{}, sdk.ErrUnknownRequest("Access code is not valid!")
    }

    ids := keeper.FindBy(ctx, path[1], path[2], path[3])

    res, err := codec.MarshalJSONIndent(keeper.cdc, ids)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}

func queryAdminGroup(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
    adminAddresses := keeper.ShowAdminGroup(ctx)

    res, err := codec.MarshalJSONIndent(keeper.cdc, adminAddresses)
    if err != nil {
        panic("could not marshal result to JSON")
    }

    return res, nil
}


//////////////////
//              //
// helper funcs //
//              //
//////////////////

func verifyAccessCode(accessCode string) (sdk.AccAddress, error) {
    r1 := regexp.MustCompile("-")
    r2 := regexp.MustCompile("_")
    accessCode1 := r1.ReplaceAllString(accessCode, "+");
    accessCode2 := r2.ReplaceAllString(accessCode1, "/");

    parts := strings.Split(accessCode2, ":")
    pubKeyBytes, _ := base64.StdEncoding.DecodeString(parts[0])
    timeStamp      := parts[1]
    signature, _   := base64.StdEncoding.DecodeString(parts[2])

    //pubKeyBytes, _ := hex.DecodeString(pubKeyStr)
    //pubKey, _ := crypto.PubKey(hex.DecodeString(pubKeyStr))

    var pubKey secp256k1.PubKeySecp256k1
    copy(pubKey[:], pubKeyBytes)
    //pubKey := crypto.PubKey(pubKeyBytes)

    if ! pubKey.VerifyBytes([]byte(timeStamp), []byte(signature)) {
        return nil, errors.New("Failed to verify signature")
    }

    timeStampInt, err := strconv.Atoi(timeStamp)
    if err != nil {
        return nil, errors.New("Failed to verify access token")
    }
    now := time.Now().UnixNano() / 1000000
    diff := now - int64(timeStampInt)
    if diff < 0 { diff -= 0 }

    if diff < MaxAllowedTimeDiff {
        address := sdk.AccAddress(pubKey.Address())
        return address, nil
    } else {
        return nil, errors.New("Failed to verify access token")
    }
}

