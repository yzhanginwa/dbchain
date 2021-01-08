package tailor_lua

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestCheckLuaLoop(t *testing.T) {
	scripts := []string{
		`a = 5 ;function double(num1);for i=10,1,-1 do;print(i);end;result = double(num1);return result;end`,
		`function double(num1);for i=10,1,-1 do;print(i);end;result = double(num1);return result;end;a = 5;`,
		`function double(num1);for i=10,1,-1 do;print(i);end;result = double(num1);return result;end`,
		`function double(num1);a = 10;repeat;print("a的值为:", a);a = a + 1;until( a > 15 );end;`,
		`a = 5`,
		"*************",
		`function hello();print("hello");end;`,
	}
	errs := []string{
		"Only one function can be defined ",
		"Only one function can be defined ",
		"can not use loop in a function",
		"can not use loop in a function",
		"Only starting with defined function is supported ",
		"<string> line:1(column:1) near '*':   syntax error\n",
		"",
	}
	for i, v := range scripts {
		err := CompileAndCheckLuaScript(v)
		if i == len(scripts)-1 {
			if err != nil {
				t.Fatal(err.Error())
			}
		} else {
			assert.Equal(t, err.Error(), errs[i], fmt.Sprintf("failed ! got : %s; want : %s", err.Error(), errs[i]))
		}
	}
}
