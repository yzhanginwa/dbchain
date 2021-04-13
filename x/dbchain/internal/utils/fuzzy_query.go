package utils

import (
	"errors"
	"regexp"
	"strings"
)

func DealFuzzyQueryString(likeExpress string )  (*regexp.Regexp, error){
	// only support \\ \% \_
	temp := strings.ReplaceAll(likeExpress,`\\`,"")
	temp = strings.ReplaceAll(temp, `\%`, "")
	temp = strings.ReplaceAll(temp, `\_`, "")
	if strings.Contains(temp, `\`) {
		return nil, errors.New("express error")
	}
	//likeExpress := `%abc\\`
	special := map[byte]string {
		'$' : `\$`,
		'(' : `\(`,
		')' : `\)`,
		'*' : `\*`,
		'+' : `\+`,
		'.' : `\.`,
		'[' : `\[`,
		']' : `\]`,
		'?' : `\?`,
		'^' : `\^`,
		'{' : `\{`,
		'}' : `\}`,
		'|' : `\|`,
	}
	//replace special char
	likeStringTemp := ""
	for _,char := range likeExpress {
		if v , ok := special[byte(char)]; ok {
			likeStringTemp += v
		} else {
			likeStringTemp += string(char)
		}
	}
	likeExpress = likeStringTemp

	//start without %, it should be add ^
	//end without %, expect \%, it should be add $
	if !strings.HasPrefix(likeExpress, `%`) {
		likeExpress = `^` + likeExpress
	}
	if !(strings.HasSuffix(likeExpress, `%`) && !strings.HasSuffix(likeExpress, `\%`)) {
		likeExpress += `$`
	}

	//deal % and _
	sSplit := strings.Split(likeExpress, `\%`)
	for i,v := range sSplit {
		sSplit[i] = strings.ReplaceAll(v,`%`, `.*`)
	}

	likeExpress = strings.Join(sSplit,`%`)
	sSplit = strings.Split(likeExpress, `\_`)
	for i,v := range sSplit {
		sSplit[i] = strings.ReplaceAll(v,`_`, `.`)
	}
	likeExpress = strings.Join(sSplit, `_`)
	reg , err := regexp.Compile(likeExpress)
	if err != nil {
		return nil, err
	}
	return reg, nil
}