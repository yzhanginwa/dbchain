package types

import (
    "strings"
    "fmt"
)


////////////////

// QueryTables Result table names
type QueryTables []string

// implement fmt.Stringer
func (t QueryTables) String() string {
    return strings.Join(t, "\n")
}

////////////////

type QueryRowFields map[string]string

func (rf QueryRowFields) String() string {
    var result = ""
    for k, v := range rf {
        result = fmt.Sprintf("%s%s: %s\n", result, k, v)
    }
    return result
}

////////////////

// QueryTables Result table names
type QuerySliceOfString []string

// implement fmt.Stringer
func (t QuerySliceOfString) String() string {
    return strings.Join(t, "\n")
}

////////////////
