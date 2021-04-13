package utils

import (
	"fmt"
	"testing"
)

func TestDealFuzzyQueryString(t *testing.T){
	resource := map[string]string {
		`%abc%` : `xxxabcxxx`,	//true
		`ab%c`  : `abxxxxc`,	//true
		`_abc_` : `xabcx`,		//true
		`a_bc`  : `axbc`,		//true
	}
	for k, v := range resource {
		reg , err := DealFuzzyQueryString(k)
		if err != nil {
			t.Error(fmt.Sprintf("%s : %s err", k, v))
		}
		if !reg.MatchString(v) {
			t.Error(fmt.Sprintf("%s : %s err", k, v))
		}
	}


	resource = map[string]string {
		`%abc%` : `xxxabyyyycxxx`,	    //false
		`ab%c`  : `yyyyabxxxxcyy`,	    //false
		`_abc_` : `yyyyxabcxyyyy`,		//false
		`a_bc`  : `yyyyaxbcyyyyy`,		//false
	}
	for k, v := range resource {
		reg , err := DealFuzzyQueryString(k)
		if err != nil {
			t.Error(fmt.Sprintf("%s : %s err", k, v))
		}
		if reg.MatchString(v) {
			t.Error(fmt.Sprintf("%s : %s err", k, v))
		}
	}

	resource = map[string]string {
		`ab\\c`: `ab\c`,	    //true
		`ab*c` : `ab*c`,	    //true
		`ab.c` : `ab.c`,		//true
		`ab?c` : `ab?c`,		//true
		`ab\%c`: `ab%c`,
		`ab\_c`: `ab_c`,
		`abc$` : `abc$`,
		`abc(` : `abc(`,
		`abc)` : `abc)`,
		`abc*` : `abc*`,
		`abc+` : `abc+`,
		`abc.` : `abc.`,
		`abc[` : `abc[`,
		`abc]` : `abc]`,
		`abc?` : `abc?`,
		`abc^` : `abc^`,
		`abc{` : `abc{`,
		`abc}` : `abc}`,
		`abc|` : `abc|`,
	}

	for k, v := range resource {
		reg , err := DealFuzzyQueryString(k)
		if err != nil {
			t.Error(fmt.Sprintf("%s : %s err", k, v))
		}
		if !reg.MatchString(v) {
			t.Error(fmt.Sprintf("%s : %s err", k, v))
		}
	}

	resource = map[string]string {
		`ab\c`  : `ab\c`,	    //false
		`ab\*c` : `ab*c`,
		`ab\.c` : `ab.c`,
		`ab\?c` : `ab?c`,
		`abc\$` : `abc$`,
		`abc\(` : `abc(`,
		`abc\)` : `abc)`,
		`abc\*` : `abc*`,
		`abc\+` : `abc+`,
		`abc\.` : `abc.`,
		`abc\[` : `abc[`,
		`abc\]` : `abc]`,
		`abc\?` : `abc?`,
		`abc\^` : `abc^`,
		`abc\{` : `abc{`,
		`abc\}` : `abc}`,
		`abc\|` : `abc|`,
	}

	for k, v := range resource {
		_ , err := DealFuzzyQueryString(k)
		if err == nil {
			t.Error(fmt.Sprintf("%s : %s err", k, v))
		}
	}
}