package keeper

import (
    "fmt"
    "strings"
    "strconv"
    "errors"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/other"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
    ss "github.com/yzhanginwa/dbchain/x/dbchain/internal/super_script"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/super_script/eval"
    shell "github.com/ipfs/go-ipfs-api"
)


func (k Keeper) Insert(ctx sdk.Context, appId uint, tableName string, fields types.RowFields, owner sdk.AccAddress) (uint, error){
    id, err := k.PreInsertCheck(ctx, appId, tableName, fields, owner)
    if err != nil {
        return id, err
    }

    // as far the first go routine to be used
    go k.tryToPinFile(ctx, appId, tableName, fields, owner)

    id, err = getNextId(k, ctx, appId, tableName)
    if err != nil {
        return 0, errors.New(fmt.Sprintf("Failed to get id for table %s", tableName))
    }

    // to set the 2 special fields
    fields["id"] = strconv.Itoa(int(id))
    fields["created_by"] = owner.String()
    fields["created_at"] = other.GetCurrentBlockTime().String()

    id, err = k.Write(ctx, appId, tableName, id, fields, owner)
    if err != nil {
        return id, err
    }

    k.updateIndex(ctx, appId, tableName, id, fields)
    k.applyTrigger(ctx, appId, tableName, fields, owner)
    return id, nil
}

// TODO: need to think over how and when to allow updating
func (k Keeper) Update(ctx sdk.Context, appId uint, tableName string, id uint, fields types.RowFields, owner sdk.AccAddress) (uint, error){
    // TODO: need to check the ownership of the record
    k.Write(ctx, appId, tableName, id, fields, owner)
    k.updateIndex(ctx, appId, tableName, id, fields)
    return id, nil
}


func (k Keeper) Write(ctx sdk.Context, appId uint, tableName string, id uint, fields types.RowFields, owner sdk.AccAddress) (uint, error){
    store := ctx.KVStore(k.storeKey)

    fieldNames, err := k.getTableFields(ctx, appId, tableName)
    if err != nil {
        return 0, errors.New(fmt.Sprintf("Failed to get fields for table %s", tableName))
    }

    if id == 0 {
        return 0, errors.New(fmt.Sprintf("Id for table %s is invalid", tableName))
    }

    for _, fieldName := range fieldNames {
        if value, ok := fields[fieldName]; ok {
            key := getDataKeyBytes(appId, tableName, fieldName, id)
            store.Set(key, k.cdc.MustMarshalBinaryBare(value))
        }
    }

    return id, nil
}

func (k Keeper) Delete(ctx sdk.Context, appId uint, tableName string, id uint, owner sdk.AccAddress) (uint, error){
    store := ctx.KVStore(k.storeKey)

    fieldNames, err := k.getTableFields(ctx, appId, tableName)
    if err != nil {
        return 0, errors.New(fmt.Sprintf("Failed to get fields for table %s", tableName))
    }

    if id == 0 {
        return 0, errors.New("Id cannot be empty")
    }

    for _, fieldName := range fieldNames {
        key := getDataKeyBytes(appId, tableName, fieldName, id)
        store.Delete(key)
    }

    // TODO: to remove the related indexes
    return id, nil
}

func (k Keeper) Freeze(ctx sdk.Context, appId uint, tableName string, id uint, owner sdk.AccAddress) (uint, error){
    store := ctx.KVStore(k.storeKey)

    if id == 0 {
        return 0, errors.New("Id cannot be empty")
    }

    keyAt := getDataKeyBytes(appId, tableName, types.FLD_FROZEN_AT, id)
    bz := store.Get(keyAt)
    if bz != nil {
        return id, nil
    }
    store.Set(keyAt, k.cdc.MustMarshalBinaryBare(other.GetCurrentBlockTime().String()))

    keyBy := getDataKeyBytes(appId, tableName, types.FLD_FROZEN_BY, id)
    store.Set(keyBy, k.cdc.MustMarshalBinaryBare(owner.String()))

    // TODO: to remove the related indexes
    return id, nil
}

//////////////////
//              //
// helper funcs //
//              //
//////////////////

func (k Keeper) PreInsertCheck(ctx sdk.Context, appId uint, tableName string, fields types.RowFields, owner sdk.AccAddress) (uint, error) {
    if !k.IsDatabaseUser(ctx, appId, owner) {
        return 0, errors.New(fmt.Sprintf("Do not have user permission on database %d", appId))
    }
    if(!k.haveWritePermission(ctx, appId, tableName, owner)) {
        return 0, errors.New(fmt.Sprintf("Do not have permission inserting table %s", tableName))
    }

    if(!k.validateInsertion(ctx, appId, tableName, fields, owner)) {
        return 0, errors.New(fmt.Sprintf("Failed validation when inserting table %s", tableName))
    }

    if(!k.preprocessPayment(ctx, appId, tableName, fields, owner)) {
        return 0, errors.New(fmt.Sprintf("Failed validation of record of payment table %s", tableName))
    }

    return 0, nil
}

func (k Keeper) haveWritePermission(ctx sdk.Context, appId uint, tableName string, owner sdk.AccAddress) bool {
    writableGroups := k.GetWritableByGroups(ctx, appId, tableName)
    if len(writableGroups) == 0 {
        return true
    }

    for _, group := range writableGroups {
        members := k.getGroupMembers(ctx, appId, group)
        if utils.AddressIncluded(members, owner) {
            return true
        }
    }
    return false
}

func (k Keeper) preprocessPayment(ctx sdk.Context, appId uint, tableName string, fields types.RowFields, owner sdk.AccAddress) bool {
    if !k.isTablePayment(ctx, appId, tableName) {
        return true
    }

    fieldNames, err := k.getTableFields(ctx, appId, tableName)
    if err != nil {
        return false
    }

    var senderRecipient map[string]sdk.AccAddress
    var amount int

    rounds := 0   // used to count the collecting of sender, recipient, and amount
    for _, fieldName := range fieldNames {
        if fieldName == "sender" || fieldName == "recipient" {
            rounds += 1
            if value, ok := fields[fieldName]; ok {
                address, err := sdk.AccAddressFromBech32(value)
                if err == nil {
                    senderRecipient[fieldName] = address
                } else {
                    return false
                }
            } else {
                return false
            }
        } else if fieldName == "amount" {
            rounds += 1
            if value, ok := fields[fieldName]; ok {
                if amount, ok = validateAmount(value); !ok {
                    return false
                }
            } else {
                return false
            }
        }
    }
    if rounds != 3 {
        return false
    }
    coin := sdk.NewCoin("dbctoken", sdk.NewInt(int64(amount)))
    err = k.CoinKeeper.SendCoins(ctx, senderRecipient["sender"], senderRecipient["recipient"], sdk.Coins{coin})
    if err != nil {
        return false
    }
    return true
}

func (k Keeper) validateInsertion(ctx sdk.Context, appId uint, tableName string, fields types.RowFields, owner sdk.AccAddress) bool {
    if ok := k.validateInsertionWithTableOptions(ctx, appId, tableName, fields, owner); !ok {
        return false
    }
    if ok := k.validateInsertionWithFieldOptions(ctx, appId, tableName, fields, owner); !ok {
        return false
    }
    if ok := k.validateInsertionWithInsertFilter(ctx, appId, tableName, fields, owner); !ok {
        return false
    }
    return true
}

func (k Keeper) validateInsertionWithTableOptions(ctx sdk.Context, appId uint, tableName string, fields types.RowFields, owner sdk.AccAddress) bool {
    options, _ := k.GetOption(ctx, appId, tableName)
    if utils.StringIncluded(options, string(types.TBLOPT_AUTH)) {
        if checkWithOracleAuth(k, ctx, fields, owner) {
            return true
        }
        return false
    }
    return true
}

func (k Keeper) validateInsertionWithFieldOptions(ctx sdk.Context, appId uint, tableName string, fields types.RowFields, owner sdk.AccAddress) bool {
    fieldNames, err := k.getTableFields(ctx, appId, tableName)
    if err != nil {
        return(false)
    }

    for _, fieldName := range fieldNames {
        if(isSystemField(fieldName)) {
            continue
        }
        fieldOptions, _ := k.GetColumnOption(ctx, appId, tableName, fieldName)

        for _, opt := range fieldOptions {
            if opt == string(types.FLDOPT_INT) {
                if value, ok := fields[fieldName]; ok {
                    if _, err := strconv.Atoi(value); err != nil {
                       return false
                    }
                }
            }

            if opt == string(types.FLDOPT_NOTNULL) {
                if value, ok := fields[fieldName]; ok {
                    if(len(value) < 1) {
                        return(false)
                    }
                } else {
                  return(false)
                }
            }

            if opt == string(types.FLDOPT_UNIQUE) {
                if value, ok := fields[fieldName]; ok {
                    if(len(value)>0) {
                        ids := k.FindBy(ctx, appId, tableName, fieldName, []string{value}, owner)
                        if len(ids) > 0 {
                            return(false)
                        }
                    }
                }
            }

            if opt == string(types.FLDOPT_OWN) {
                if tn, ok := utils.GetTableNameFromForeignKey(fieldName); ok {
                    foreignId := fields[fieldName]
                    u64, err := strconv.ParseUint(foreignId, 10, 64)
                    if err != nil {
                        return false
                    }
                    foreignOwner, err := k.FindField(ctx, appId, tn, uint(u64), "created_by")
                    if err == nil {
                        if foreignOwner != owner.String() {
                            return false
                        }
                    } else {
                        return false
                    }
                } else {
                    return false
                }
            }

            if types.ValidateEnumColumnOption(opt) {
                if value, ok := fields[fieldName]; ok {
                    items := types.GetEnumColumnOptionItems(opt)
                    if utils.StringIncluded(items, value) {
                        return true
                    }
                    return false
                }
            }
        }
    }
    return true
}

func (k Keeper) applyTrigger(ctx sdk.Context, appId uint, tableName string, fields types.RowFields, owner sdk.AccAddress) {
    table, err := k.GetTable(ctx, appId, tableName)
    if err != nil {
        return
    }

    trigger := table.Trigger
    if len(trigger) == 0 {
        return
    }

    k.runFilterOrTrigger(ctx, appId, tableName, fields, owner, ss.TRIGGER, trigger)
}

func (k Keeper) validateInsertionWithInsertFilter(ctx sdk.Context, appId uint, tableName string, fields types.RowFields, owner sdk.AccAddress) bool {
    table, err := k.GetTable(ctx, appId, tableName)
    if err != nil {
        return false
    }

    filter := table.Filter
    if len(filter) == 0 {
        return true
    }

    result := k.runFilterOrTrigger(ctx, appId, tableName, fields, owner, ss.FILTER, filter)
    return result
}

func (k Keeper) runFilterOrTrigger(ctx sdk.Context, appId uint, tableName string, fields types.RowFields, owner sdk.AccAddress, scriptType ss.ScriptType, script string) bool {
    //TODO: create a database associated cache mapping of table script to syntax tree
    //so that we don't have to parse the script for each insertion

    fn1 := getScriptValidationCallbackOne(k, ctx, appId, tableName)
    fn2 := getScriptValidationCallbackTwo(k, ctx, appId, tableName)

    parser := ss.NewParser(strings.NewReader(script), fn1, fn2)
    var err error
    if scriptType == ss.FILTER {
        err = parser.ParseFilter()
    } else {
        err = parser.ParseTrigger()
    }

    if err != nil {
        return false
    }

    fieldValueCallback := getGetFieldValueCallback(k, ctx, appId, owner)
    tableValueCallback := getGetTableValueCallback(k, ctx, appId, owner)
    insertCallback     := getInsertCallback(k, ctx, appId, owner)

    program := eval.NewProgram(tableName, fields, script, fieldValueCallback, tableValueCallback, insertCallback)
    result := program.EvaluateScript(parser.GetSyntaxTree())

    return result
}

func (k Keeper) tryToPinFile(ctx sdk.Context, appId uint, tableName string, fields types.RowFields, owner sdk.AccAddress) bool {
    fieldNames, err := k.getTableFields(ctx, appId, tableName)
    if err != nil {
        return(false)
    }

    for _, fieldName := range fieldNames {
        if(isSystemField(fieldName)) {
            continue
        }
        fieldOptions, _ := k.GetColumnOption(ctx, appId, tableName, fieldName)
        if(utils.ItemExists(fieldOptions, string(types.FLDOPT_FILE))) {
            if value, ok := fields[fieldName]; ok {
                sh := shell.NewShell("localhost:5001")
                err =sh.Pin(value)
                if err != nil {
                    logger := k.Logger(ctx)
                    logger.Error(fmt.Sprintf("Failed to pin ipfs cid %s", value))
                    return false
                }
            }
        }
    }
    return true
}

func (k Keeper) isTablePayment(ctx sdk.Context, appId uint, tableName string) bool {
    tableOptions, _ := k.GetOption(ctx, appId, tableName)
    return utils.ItemExists(tableOptions, string(types.TBLOPT_PAYMENT))
}

func validateAmount(amount string) (int, bool) {
    i, err := strconv.Atoi(amount)
    if err != nil {
        return 0, false
    }

    if i > 0 {
        return i, true
    } else {
        return 0, false
    }
}
