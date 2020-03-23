package types

import (
    "strings"
    "fmt"
    sdk "github.com/cosmos/cosmos-sdk/types"
)

/////////////////
//             //
// QueryTables //
//             //
/////////////////

// QueryTables Result table names
type QueryTables []string

// implement fmt.Stringer
func (t QueryTables) String() string {
    return strings.Join(t, "\n")
}

////////////////////
//                //
// QueryRowFields //
//                //
////////////////////

type QueryRowFields map[string]string

func (rf QueryRowFields) String() string {
    var result = ""
    for k, v := range rf {
        result = fmt.Sprintf("%s%s: %s\n", result, k, v)
    }
    return result
}

////////////////////////
//                    //
// QuerySliceOfString //
//                    //
////////////////////////

type QuerySliceOfString []string

// implement fmt.Stringer
func (t QuerySliceOfString) String() string {
    return strings.Join(t, "\n")
}

///////////////////
//               //
// QueryOfString //
//               //
///////////////////

type QueryOfString string

// implement fmt.Stringer
func (t QueryOfString) String() string {
    return string(t)
}

////////////////////
//                //
// QueryOfBoolean //
//                //
////////////////////

type QueryOfBoolean bool

// implement fmt.Stringer
func (t QueryOfBoolean) String() string {
    if t {
        return "true"
    } else {
        return "false"
    }
}

/////////////////////
//                 //
// QueryAdminGroup //
//                 //
/////////////////////

type QueryAdminGroup []sdk.AccAddress

func (ag QueryAdminGroup) String() string {
    var buf []string
    for index, addr := range ag {
        buf[index] = fmt.Sprintf("%s", addr)
    }
    return strings.Join(buf, "\n")
}
