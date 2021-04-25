package types

import (
    "fmt"
    "strings"
    sdk "github.com/cosmos/cosmos-sdk/types"
)

type RowFields map[string]string
type RowFieldsJson []byte

/////////////////
//             //
// application //
//             //
/////////////////

type Database struct {
    Owner sdk.AccAddress    `json:"owner"`
    AppCode string          `json:"appcode"`
    AppId uint              `json:"appid"`
    Name string             `json:"name"`
    Description string      `json:"description"`
    PermissionRequired bool `json:"permission_required"`
    Deleted bool            `json:"deleted"`
    Expiration int64        `json:"expiration"`
    SchemaFrozen bool       `json:"schema_frozen"`
    DataFrozen bool         `json:"data_frozen"`
}

func NewDatabase() Database {
    return Database{
        SchemaFrozen: false,
        DataFrozen: false,
    }
}

func (d Database) String() string {
    return strings.TrimSpace(fmt.Sprintf(`AppCoe: %s`, d.AppCode))
}

//////////////
//          //
// function //
//          //
//////////////

type Function struct {
    Owner sdk.AccAddress   `json:"owner"`
    Name string            `json:"name"`
    Description string     `json:"description"`
    Body string            `json:"body"`
}

func NewFunction() Function {
    return Function{}
}

////////////
//        //
// friend //
//        //
////////////

type Friend struct {
    Address string `json:"address"`
    Name    string `json:"name"`
}

func NewFriend() Friend {
    return Friend{}
}

func (f Friend) String() string {
    return strings.TrimSpace(fmt.Sprintf(`Addr: %s`, f.Address))
}

///////////
//       //
// table //
//       //
///////////

// Memo is for the table and Memos are for the fields
type Table struct {
    Owner sdk.AccAddress      `json:"owner"`
    Name string               `json:"name"`
    Fields []string           `json:"fields"`
    Memos []string            `json:"memos"`
    Filter string             `json:"filter"`
    Trigger string            `json:"trigger"`
    Memo string               `json:"memo"`
}

func NewTable() Table {
    return Table {}
}

// implement fmt.Stringer
func (t Table) String() string {
    return strings.TrimSpace(fmt.Sprintf(`Name: %s`, t.Name))
}

//////////////////////////
//                      //
// tx status            //
//                      //
//////////////////////////
type  TxStatus struct {
    State string  `json:"state"`//success,fail,processing
    Index int    `json:"index"`//count of message
    Err  string    `json:"err"`
    timeStamp int64
}

func NewTxStatus(state string, index int, err string, ts int64) *TxStatus {
    return &TxStatus {
        State: state,
        Index: index,
        Err: err,
        timeStamp: ts,
    }
}
func (t *TxStatus)GetTimeStamp() int64 {
    return t.timeStamp
}

//////////////////////////
//                      //
// reserved field names //
//                      //
//////////////////////////

const (
    FLD_FROZEN_AT       string = "_frozen_at_"
    FLD_FROZEN_BY       string = "_frozen_by_"
)

//////////////////
//              //
// option types //
//              //
//////////////////

type TableOption string

const (
    TBLOPT_PUBLIC      TableOption = "public"
    TBLOPT_WRITABLE_BY TableOption = "writable-by"  // writable-by: the table can only be written by members of writable-by group
    TBLOPT_PAYMENT     TableOption = "payment"      // payment: this table needs to have fields "sender", "recipient", token_name, and "amount".
                                                    //          after the a row is saved, the amount of token_name is sent from sender to recipient
    TBLOPT_UPDATABLE   TableOption = "updatable"
    TBLOPT_DELETABLE   TableOption = "deletable"
    TBLOPT_AUTH        TableOption = "auth"
)

type FieldOption string

const (
    FLDOPT_NOTNULL    FieldOption = "not-null"
    FLDOPT_UNIQUE     FieldOption = "unique"
    FLDOPT_OWN        FieldOption = "own"
)

type FieldDataType string

const (
    FLDTYP_STRING     FieldDataType = "string"
    FLDTYP_FILE       FieldDataType = "file"
    FLDTYP_INT        FieldDataType = "int"
    FLDTYP_DECIMAL    FieldDataType = "decimal"
)
