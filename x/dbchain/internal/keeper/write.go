package keeper

import (
    "errors"
    "fmt"
    sdk "github.com/dbchaincloud/cosmos-sdk/types"
    shell "github.com/ipfs/go-ipfs-api"
    lua "github.com/yuin/gopher-lua"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/other"
    ss "github.com/yzhanginwa/dbchain/x/dbchain/internal/super_script"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/super_script/eval"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/types"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/utils"
    "strconv"
    "strings"
    storeTypes "github.com/dbchaincloud/cosmos-sdk/store/types"
)

const fileSizeGasRate = 10
func (k Keeper) Insert(ctx sdk.Context, appId uint, tableName string, fields types.RowFields, owner sdk.AccAddress) (uint, error){
    return k.InsertCore(ctx, appId, tableName, fields, owner, true)
}

func (k Keeper) InsertCore(ctx sdk.Context, appId uint, tableName string, fields types.RowFields, owner sdk.AccAddress, IsCallTrigger bool ) (uint, error) {
    L := lua.NewState(lua.Options{
        SkipOpenLibs : true,
        RegistrySize: 32,
    })
    L.SetGlobal("IsRegisterData",lua.LBool(false))
    openBase(L)
    registerTableType(L, ctx, appId, k, owner)

    defer L.Close()
    id, err := k.PreInsertCheck(ctx, appId, tableName, fields, owner, L)
    if err != nil {
        return id, err
    }

    // as far the first go routine to be used
    _, allUploadFileSize := k.tryToPinFile(ctx, appId, tableName, fields, owner)
    id, err = getNextId(k, ctx, appId, tableName)
    if err != nil {
        return 0, errors.New(fmt.Sprintf("Failed to get id for table %s", tableName))
    }

    // to set the 2 special fields
    fields["id"] = strconv.Itoa(int(id))
    fields["created_by"] = owner.String()
    //测试改为时间戳
    fields["created_at"] = fmt.Sprintf("%d",other.GetCurrentBlockTime().UnixNano()/(1000*1000))
    fields["tx_hash"] = k.GetTxHash(ctx)

    id, err = k.Write(ctx, appId, tableName, id, fields, owner)
    if err != nil {
        return id, err
    }

    _, err = k.appendIndexForRow(ctx, appId, tableName, id)
    if err != nil {
        return id, err
    }

    if IsCallTrigger {
        k.applyTrigger(ctx, appId, tableName, fields, owner, L)
    }
    k.consumeGasByUploadFile(ctx, allUploadFileSize)
    if ctx.GasMeter().IsOutOfGas() {
        return 0, errors.New("out of gas")
    }
    return id, nil
}
// TODO: need to think over how and when to allow updating
func (k Keeper) Update(ctx sdk.Context, appId uint, tableName string, id uint, fields types.RowFields, owner sdk.AccAddress) (uint, error){
//    // TODO: need to check the ownership of the record
//    k.Write(ctx, appId, tableName, id, fields, owner)
//    k.updateIndex(ctx, appId, tableName, id, fields)
    return id, nil
}


func (k Keeper) Write(ctx sdk.Context, appId uint, tableName string, id uint, fields types.RowFields, owner sdk.AccAddress) (uint, error){
    store := DbChainStore(ctx, k.storeKey)

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
            err := store.Set(key, k.cdc.MustMarshalBinaryBare(value))
            if err != nil{
                return 0,err
            }
        }
    }

    return id, nil
}

func (k Keeper) Delete(ctx sdk.Context, appId uint, tableName string, id uint, owner sdk.AccAddress) (uint, error){
    store := DbChainStore(ctx, k.storeKey)

    fieldNames, err := k.getTableFields(ctx, appId, tableName)
    if err != nil {
        return 0, errors.New(fmt.Sprintf("Failed to get fields for table %s", tableName))
    }

    if id == 0 {
        return 0, errors.New("Id cannot be empty")
    }

    if !k.isOwnId(ctx, appId, tableName, id, owner) {
        return 0, errors.New("no permission")
    }
    cids := make([]string, 0)
    if !isRowFrozen(store, appId, tableName, id) {
        cids = k.getCids(ctx, appId, tableName, id)
    }
    for _, fieldName := range fieldNames {
        key := getDataKeyBytes(appId, tableName, fieldName, id)
        err := store.Delete(key)
        if err != nil{
            return 0, err
        }
    }

    k.RestoreVolume(ctx, appId, cids, owner.String())

    _, err = k.dropIndexForRow(ctx, appId, tableName, id)
    if err != nil {
        return 0, err
    }
    return id, nil
}

func (k Keeper) Freeze(ctx sdk.Context, appId uint, tableName string, id uint, owner sdk.AccAddress) (uint, error){
    store := DbChainStore(ctx, k.storeKey)

    if id == 0 {
        return 0, errors.New("Id cannot be empty")
    }

    if !k.isOwnId(ctx, appId, tableName, id, owner) {
        return 0, errors.New("no permission")
    }
    keyAt := getDataKeyBytes(appId, tableName, types.FLD_FROZEN_AT, id)
    bz, err := store.Get(keyAt)
    if err != nil{
        return 0, err
    }
    if bz != nil {
        return id, errors.New("Record is already frozen")
    }
    cids := k.getCids(ctx, appId, tableName, id)
    store.Set(keyAt, k.cdc.MustMarshalBinaryBare(other.GetCurrentBlockTime().String()))

    keyBy := getDataKeyBytes(appId, tableName, types.FLD_FROZEN_BY, id)
    store.Set(keyBy, k.cdc.MustMarshalBinaryBare(owner.String()))
    if len(cids) > 0 {
        k.RestoreVolume(ctx, appId, cids, owner.String())
    }
    _, err = k.dropIndexForRow(ctx, appId, tableName, id)
    if err != nil {
        return 0, err
    }
    return id, nil
}

//////////////////
//              //
// helper funcs //
//              //
//////////////////

func (k Keeper) PreInsertCheck(ctx sdk.Context, appId uint, tableName string, fields types.RowFields, owner sdk.AccAddress, L *lua.LState) (uint, error) {
    if !k.IsDatabaseUser(ctx, appId, owner) {
        return 0, errors.New(fmt.Sprintf("Do not have user permission on database %d", appId))
    }
    if(!k.haveWritePermission(ctx, appId, tableName, owner)) {
        return 0, errors.New(fmt.Sprintf("Do not have permission inserting table %s", tableName))
    }

    if(!k.validateInsertion(ctx, appId, tableName, fields, owner, L)) {
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

    var senderRecipient = make(map[string]sdk.AccAddress,0)
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

func (k Keeper) validateInsertion(ctx sdk.Context, appId uint, tableName string, fields types.RowFields, owner sdk.AccAddress, L *lua.LState) bool {
    if ok := k.validateInsertionWithTableOptions(ctx, appId, tableName, fields, owner); !ok {
        return false
    }
    if ok := k.validateInsertionWithFieldOptions(ctx, appId, tableName, fields, owner); !ok {
        return false
    }
    if ok := k.validateInsertionWithFieldDataType(ctx, appId, tableName, fields, owner); !ok {
        return false
    }
    if ok := k.validateInsertionWithInsertFilter(ctx, appId, tableName, fields, owner, L); !ok {
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

func (k Keeper) validateInsertionWithFieldDataType(ctx sdk.Context, appId uint, tableName string, fields types.RowFields, owner sdk.AccAddress) bool {
    fieldNames, err := k.getTableFields(ctx, appId, tableName)
    if err != nil {
        return(false)
    }

    for _, fieldName := range fieldNames {
        if(isSystemField(fieldName)) {
            continue
        }
        fieldDataType, _ := k.GetColumnDataType(ctx, appId, tableName, fieldName)

        if fieldDataType == string(types.FLDTYP_INT) {
            if value, ok := fields[fieldName]; ok {
                if _, err := strconv.Atoi(value); err != nil {
                    return false
                }
            }
        } else if fieldDataType == string(types.FLDTYP_DECIMAL) {
            if value, ok := fields[fieldName]; ok {
                if _, err := strconv.ParseFloat(value,64); err != nil {
                    return false
                }
            }
        } else if fieldDataType == string(types.FLDTYP_ADDRESS) {
            if value, ok := fields[fieldName]; ok {
                if _, err := sdk.AccAddressFromBech32(value); err != nil {
                    return false
                }
            }
        } else if fieldDataType ==  string(types.FLDTYP_TIME) {
            if value, ok := fields[fieldName]; ok {
                num , err := strconv.ParseInt(value, 10, 64)
                if err != nil || num < 0 {
                    return false
                }
            }
        }
    }
    return true
}

func (k Keeper) applyTrigger(ctx sdk.Context, appId uint, tableName string, fields types.RowFields, owner sdk.AccAddress, L *lua.LState) {
    table, err := k.GetTable(ctx, appId, tableName)
    if err != nil {
        return
    }

    trigger := table.Trigger
    if len(trigger) == 0 {
        return
    }

    k.runLuaFilter(ctx, appId, tableName, fields, owner, trigger, L)
}

func (k Keeper) validateInsertionWithInsertFilter(ctx sdk.Context, appId uint, tableName string, fields types.RowFields, owner sdk.AccAddress, L *lua.LState) bool {
    table, err := k.GetTable(ctx, appId, tableName)
    if err != nil {
        return false
    }

    filter := table.Filter
    if len(filter) == 0 {
        return true
    }

    result := k.runLuaFilter(ctx, appId, tableName, fields, owner, filter, L)
    return result
}

func (k Keeper) runLuaFilter(ctx sdk.Context, appId uint, tableName string, fields types.RowFields, owner sdk.AccAddress, script string, L *lua.LState) bool {
    if L.GetGlobal("IsRegisterData") == lua.LFalse {
        if !k.registerThisData(ctx, appId, tableName, fields, owner, L, script) {
            return false
        }
        L.SetGlobal("IsRegisterData",lua.LBool(true))
        //point : get go function
        goExportFunc := getGoExportFilterFunc(ctx, appId, k, owner)
        //register go function
        for name, fn := range goExportFunc{
            L.SetGlobal(name, L.NewFunction(fn))
        }
    }
    newScript := restructureLuaScript(script)
    if err := L.DoString(newScript); err != nil{
        return false
    }

    //handle return
    strSuccess := L.Get(1).String()
    defer L.Pop(L.GetTop())
    if strSuccess == "true" {
        return true
    }
    return false
}

func (k Keeper) registerThisData(ctx sdk.Context, appId uint, tableName string, fields types.RowFields, owner sdk.AccAddress, L *lua.LState, script string) bool{
    this := L.NewTable()
    tbFields, err := k.getTableFields(ctx, appId, tableName)
    if err != nil { return false }
    for _, field := range tbFields {
        isForeignKey := false
        v , ok := fields[field]
        if ok {
            if strings.Contains(field,"_") {
                temp := "." + field + ".parent"
                if strings.Contains(script,temp) {
                    isForeignKey = true
                }
            }
            if !isForeignKey {
                dataType, err := k.GetColumnDataType(ctx, appId, tableName, field)
                if err == nil && (dataType == string(types.FLDTYP_INT) || dataType == string(types.FLDTYP_DECIMAL)){
                    number , err := strconv.ParseFloat(v,64)
                    if err == nil {
                        this.RawSetString(field, lua.LNumber(number))
                    } else {
                        this.RawSetString(field, lua.LString(v))
                    }

                } else {
                    this.RawSetString(field, lua.LString(v))
                }

            }
        } else if field == "created_by"{
            this.RawSetString(field, lua.LString(owner.String()))
        } else {
            this.RawSetString(field, lua.LString("nil"))
        }

        if isForeignKey {
            foreignKey := L.NewTable()
            parent := L.NewTable()
            tableAndKey := strings.Split(field,"_")
            if len(tableAndKey) != 2 {
                return false
            }

            parentTableName := tableAndKey[0]
            parentTableField := tableAndKey[1]

            ids := k.FindBy(ctx, appId, parentTableName, parentTableField, []string{v}, owner)
            if len(ids) != 1 {
                return false
            }
            RowFields, err := k.DoFind(ctx, appId, parentTableName, ids[0])
            if err != nil {
                return false
            }

            for key , value  := range RowFields {
                parent.RawSetString(key, lua.LString(value))
            }
            foreignKey.RawSetString("parent", parent)
            this.RawSetString(field, foreignKey)
        }
    }
    L.SetGlobal("this", this) //register this table
    return true
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

func (k Keeper) tryToPinFile(ctx sdk.Context, appId uint, tableName string, fields types.RowFields, owner sdk.AccAddress) (result bool, allUploadFileSize uint64) {

    fieldNames, err := k.getTableFields(ctx, appId, tableName)
    if err != nil {
        return(false), 0
    }

    for _, fieldName := range fieldNames {
        if(isSystemField(fieldName)) {
            continue
        }
        fieldDataType, _ := k.GetColumnDataType(ctx, appId, tableName, fieldName)
        if fieldDataType == string(types.FLDTYP_FILE) {
            if value, ok := fields[fieldName]; ok {
                sh := shell.NewShell("localhost:5001")
                size := getUploadFileSize(sh, value)
                allUploadFileSize += size
                k.UpdateAppUserUsedFileVolume(ctx, appId, owner.String(), fmt.Sprintf("%d", size))
                go func(sh *shell.Shell, value string) {
                    err =sh.Pin(value)
                    if err != nil {
                        logger := k.Logger(ctx)
                        logger.Error(fmt.Sprintf("Failed to pin ipfs cid %s", value))
                    }
                }(sh, value)
            }
        }
    }
    return true, allUploadFileSize
}

func (k Keeper) consumeGasByUploadFile(ctx sdk.Context, fileSize uint64)  {
    if fileSize <= 0 {
        return
    }
    gas := storeTypes.Gas(int64(fileSize)) / fileSizeGasRate
    ctx.GasMeter().ConsumeGas(gas,"consume gas by upload file")
    return
}
func (k Keeper) isTablePayment(ctx sdk.Context, appId uint, tableName string) bool {
    tableOptions, _ := k.GetOption(ctx, appId, tableName)
    return utils.ItemExists(tableOptions, string(types.TBLOPT_PAYMENT))
}

func (k Keeper) isOwnId(ctx sdk.Context, appId uint, tableName string, id uint, user sdk.AccAddress) bool {
    res , err := k.Find(ctx, appId, tableName, id, user)
    if err != nil {
        return false
    }

    if res["created_by"] == user.String() {
        return true
    }
    return  false
}

func (k Keeper) getCids(ctx sdk.Context, appId uint, tableName string, id uint) []string {
    cids := make([]string, 0)
    temp := make(map[string]bool, 0)
    fields, err := k.DoFind(ctx, appId, tableName, id)
    if err != nil {
        return cids
    }
    for field, val := range fields {
        fieldType, err := k.GetColumnDataType(ctx, appId, tableName, field)
        if err != nil {
            continue
        }
        if fieldType == string(types.FLDTYP_FILE) && !temp[fieldType]{
            temp[fieldType] = true
            cids = append(cids, val)
        }
    }
    return  cids
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

func getUploadFileSize(sh *shell.Shell, cid string) uint64{
    obj, err := sh.FileList(fmt.Sprintf("/ipfs/%s", cid))
    if err != nil {
        return 0
    }
    return obj.Size
}