package types

import (
    "regexp"
)

var (
    metaNamePattern = regexp.MustCompile("^[a-z][a-z0-9_]+$")
)

func validateMetaName(name string) bool {
    return metaNamePattern.MatchString(name)
}
