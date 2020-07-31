package types

import (
    "regexp"
)

var (
    metaNamePattern = regexp.MustCompile("^[a-z][a-z0-9_]+$")
    columnOptionEnumPattern = regexp.MustCompile(`^enum\("[^"]+"(\s*,\s*"[^"]+")*\)$`)
)

func validateMetaName(name string) bool {
    return metaNamePattern.MatchString(name)
}

func ValidateEnumColumnOption(fieldOption string) bool {
    return columnOptionEnumPattern.MatchString(fieldOption)
}
