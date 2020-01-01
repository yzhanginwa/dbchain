package types

import (
    "strings"
)


////////////////

// QueryTables Result table names
type QueryTables []string

// implement fmt.Stringer
func (t QueryTables) String() string {
    return strings.Join(t, "\n")
}

////////////////

