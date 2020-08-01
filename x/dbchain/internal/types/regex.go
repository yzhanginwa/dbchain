package types

import (
    "regexp"
)

var (
    metaNamePattern = regexp.MustCompile("^[a-z][a-z0-9_]+$")

    // enum("foo", "bar")
    columnOptionEnumPattern = regexp.MustCompile(`^enum\(("[a-z0-9_\-]+"(\s*,\s*"[a-z0-9_\-]+")*)\)$`)
    columnOptionEnumItemsSplitPattern = regexp.MustCompile(`\s*,\s*`)
)

func validateMetaName(name string) bool {
    return metaNamePattern.MatchString(name)
}

func ValidateEnumColumnOption(fieldOption string) bool {
    return columnOptionEnumPattern.MatchString(fieldOption)
}

func GetEnumColumnOptionItems(fieldOption string) []string {
    matched := columnOptionEnumPattern.FindSubmatch([]byte(fieldOption))
    if len(matched)< 2 {
        return []string{}
    }
    itemsString := string(matched[1])
    items := columnOptionEnumItemsSplitPattern.Split(itemsString, -1)

    var result []string
    for _, item := range items {
        l := len(item)
        if l < 3 {
            return []string{}
        }
        result = append(result, item[1:l-1])
    }
    return result
}
