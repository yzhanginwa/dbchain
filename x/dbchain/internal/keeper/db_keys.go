package keeper

import (
    dbk "github.com/yzhanginwa/dbchain/x/dbchain/internal/keeper/db_key"
)

//////////////////////////////////
//                              //
// application/database related //
//                              //
//////////////////////////////////

var (
    getDatabaseKey = dbk.GetDatabaseKey

    getDatabaseNextIdKey = dbk.GetDatabaseNextIdKey
    
    getDatabaseUserKey = dbk.GetDatabaseUserKey
    
    getDatabaseIteratorStartAndEndKey = dbk.GetDatabaseIteratorStartAndEndKey
    
    getAppCodeFromDatabaseKey = dbk.GetAppCodeFromDatabaseKey
    
    getDatabaseUserIteratorStartAndEndKey = dbk.GetDatabaseUserIteratorStartAndEndKey
    
    getUserFromDatabaseUserKey = dbk.GetUserFromDatabaseUserKey
    
    ///////////////////
    //               //
    // table related //
    //               //
    ///////////////////
    
    getTablesKey = dbk.GetTablesKey
    
    getNextIdKey = dbk.GetNextIdKey
    
    getTableKey = dbk.GetTableKey
    
    getMetaTableIndexKey = dbk.GetMetaTableIndexKey
    
    getTableOptionsKey = dbk.GetTableOptionsKey
    
    getColumnOptionsKey = dbk.GetColumnOptionsKey
    
    //////////////////////
    //                  //
    // function related //
    //                  //
    //////////////////////

    getFunctionKey = dbk.GetFunctionKey
    getFunctionsKey = dbk.GetFunctionsKey
    
    //////////////////
    //              //
    // data related //
    //              //
    //////////////////
    
    getIndexKey = dbk.GetIndexKey
    
    getIndexDataIteratorStartAndEndKey = dbk.GetIndexDataIteratorStartAndEndKey
    
    getDataKeyBytes = dbk.GetDataKeyBytes
    
    getFieldDataIteratorStartAndEndKey = dbk.GetFieldDataIteratorStartAndEndKey
    
    getIdFromDataKey = dbk.GetIdFromDataKey
    
    ////////////////////
    //                //
    // friend related //
    //                //
    ////////////////////
    
    getFriendKey = dbk.GetFriendKey
    
    getFriendIteratorStartAndEndKey = dbk.GetFriendIteratorStartAndEndKey
    
    getPendingFriendKey = dbk.GetPendingFriendKey
    
    getPendingFriendIteratorStartAndEndKey = dbk.GetPendingFriendIteratorStartAndEndKey
    
    ///////////////////
    //               //
    // group related //
    //               //
    ///////////////////
    
    getGroupsKey = dbk.GetGroupsKey
    
    getGroupKey = dbk.GetGroupKey
    
    getGroupMemoKey = dbk.GetGroupMemoKey
    
    getAdminGroupKey = dbk.GetAdminGroupKey
    
    //////////////////
    //              //
    // system level //
    //              //
    //////////////////
    
    getSysGroupKey = dbk.GetSysGroupKey
    
    getSysAdminGroupKey = dbk.GetSysAdminGroupKey
)
