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
    Owner sdk.AccAddress  `json:"owner"`
    AppCode string        `json:"appcode"`
    AppId uint            `json:"appid"`
    Name string           `json:"name"`
    Description string    `json:"description"`
}

func NewDatabase() Database {
    return Database{}
}

func (d Database) String() string {
    return strings.TrimSpace(fmt.Sprintf(`AppCoe: %s`, d.AppCode))
}

///////////
//       //
// table //
//       //
///////////

// the key would be like "poll:[name]"
type Table struct {
    Owner sdk.AccAddress      `json:"owner"`
    Name string               `json:"name"`
    Fields []string           `json:"fields"`
}

func NewTable() Table {
    return Table {}
}

// implement fmt.Stringer
func (t Table) String() string {
    return strings.TrimSpace(fmt.Sprintf(`Name: %s`, t.Name))
}

//////////////////
//              //
// option types //
//              //
//////////////////

type TableOption string

const (
    TBLOPT_PUBLIC     TableOption = "public"
    TBLOPT_UPDATABLE  TableOption = "updatable"
    TBLOPT_DELETABLE  TableOption = "deletable"
)

type FieldOption string

const (
    FLDOPT_NOTNULL    FieldOption = "not-null"
    FLDOPT_UNIQUE     FieldOption = "unique"
)


