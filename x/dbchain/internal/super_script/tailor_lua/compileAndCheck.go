//This package is used to precompile lua script.
//The ability of lua script is limited to  ensure the safety of dbchain system

package tailor_lua

import (
	"errors"
	"github.com/yuin/gopher-lua/ast"
	"github.com/yuin/gopher-lua/parse"
	"strings"
)

//
func CompileAndCheckLuaScript(luaScript string) error {
	chunk, err := parse.Parse(strings.NewReader(luaScript), "<string>")
	if err != nil {
		return err
	}
	if len(chunk) != 1 {
		return errors.New("Only one function can be defined ")
	}
	if _, ok := chunk[0].(*ast.FuncDefStmt); !ok {
		return errors.New("Only starting with defined function is supported ")
	} else {
		hasLoop, err := CheckLuaLoop(chunk)
		if err != nil {
			return err
		}
		if hasLoop {
			return errors.New("can not use loop in a function")
		}
	}
	return nil
}
func CheckLuaLoop(chunk []ast.Stmt) (hasLoop bool, err error) {
	defer func() {
		if rcv := recover(); rcv != nil {
			if _, ok := rcv.(*CompileError); ok {
				err = rcv.(error)
			} else {
				panic(rcv)
			}
		}
	}()
	err = nil
	context := newFuncContext("<string>", nil)
	hasLoop = compileChunk(context, chunk)
	return
}
