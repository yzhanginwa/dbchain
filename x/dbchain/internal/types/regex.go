package types

import (
    "regexp"
)

var (
    metaNamePattern = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_-]+$")
)

func validateMetaName(name string) bool {
    return metaNamePattern.MatchString(name)
}
