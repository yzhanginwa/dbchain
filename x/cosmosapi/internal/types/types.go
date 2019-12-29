package types

import (
	"fmt"
	"strings"
//	"crypto/sha256"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	KeyPrefixMeta  = "mt"
	KeyPrefixData  = "dt"
	KeyPrefixIndex = "ix"
)


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

func TableKey(name string) string {
	return fmt.Sprintf("%s:tn:%s", KeyPrefixMeta, name)
}

func NewTable() Table {
	return Table {}
}

// implement fmt.Stringer
func (t Table) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Name: %s`, t.Name))
}

