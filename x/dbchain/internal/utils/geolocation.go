package utils

import (
	"strings"
	"regexp"
)

var (
	latitudePattern  = regexp.MustCompile("^(\\+|-)?(?:90(?:(?:\\.0{1,6})?)|(?:[0-9]|[1-8][0-9])(?:(?:\\.[0-9]{1,6})?))$")
	longitudePattern = regexp.MustCompile("^(\\+|-)?(?:180(?:(?:\\.0{1,6})?)|(?:[0-9]|[1-9][0-9]|1[0-7][0-9])(?:(?:\\.[0-9]{1,6})?))$")
)

func ValidateGeolocationValue(location string) bool {
	ll := strings.Split(location, ",")
	if len(ll) != 2 {
		return false
	}

	if latitudePattern.MatchString(ll[0]) {
		if longitudePattern.MatchString(ll[1]) {
			return true
		}
	}
	return false
}
