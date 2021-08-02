package keeper

import "fmt"


const (
	KeyPrefixData  = "dt"
	KeyPrefixBsn   = "bsn"
)

func GetP2PTransferLimit() string {
	return fmt.Sprintf("%s:%s:limit", KeyPrefixBsn, KeyPrefixData)
}

func GetTokenKeeperKey() string {
	return fmt.Sprintf("%s:%s:tokenKeeper", KeyPrefixBsn, KeyPrefixData)
}


