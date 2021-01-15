package super_script

import (
	"errors"
	"fmt"
	"io"
)

//this file is built to preprocess super_script.
//Its function is to translate the super_script into Lua.

const (
	FunFieldIn = `fieldIn(`
	strComma   = `,`
	strDEQUAL  = `,"==",`
	strThen    = ` then `
	strEnd     = " end "
	strRCB     = "}"
)

type tokenStack struct {
	buf    []string
	length int
}

func newTokenStack() *tokenStack {
	return &tokenStack{
		buf:    make([]string, 0),
		length: 0,
	}
}

func (ts *tokenStack) Push(src string) {
	ts.buf = append(ts.buf, src)
	ts.length++
	return
}

func (ts *tokenStack) PushN(src []string) {
	ts.buf = append(ts.buf, src...)
	ts.length += len(src)
}

func (ts *tokenStack) Pop() (string, error) {
	if ts.length == 0 {
		return "", errors.New("stack is empty")
	}
	item := ts.buf[ts.length-1]
	ts.buf = ts.buf[:ts.length-1]
	ts.length--
	return item, nil
}

func (ts *tokenStack) PopN(n int) ([]string, error) {
	if n > ts.length {
		return nil, errors.New("stack has no so much nums")
	}
	items := make([]string, 0)
	temp := ts.buf[ts.length-n:]
	items = append(items, temp...)

	ts.buf = ts.buf[:ts.length-n]
	ts.length -= n
	return items, nil
}

func (ts *tokenStack) Clear() {
	ts.length = 0
	ts.buf = make([]string, 0)
}

type Preprocessor struct {
	ts      *tokenStack
	s       *Scanner
	tok     Token  // last read token
	lit     string // last read literal
	temp    *tokenStack
	Err     error
	Success bool
}

func NewPreprocessor(r io.Reader) *Preprocessor {
	return &Preprocessor{
		s:    NewScanner(r),
		ts:   newTokenStack(),
		temp: newTokenStack(),
	}
}

func (pc *Preprocessor) Process() {
	for {
		tok, lit := pc.s.Scan()
		if tok == EOF {
			break
		}
		pc.ts.Push(lit)
		if tok == IF || tok == ELSEIF{
			pc.tok, pc.lit = tok, lit
			if !pc.IfCondition() {
				pc.ts.Clear()
				pc.Success = false
				return
			}
		}
	}
	pc.Success = true
	return
}

func (pc *Preprocessor) Reconstruct() string {
	s := ""
	for {
		temp, err := pc.ts.Pop()
		if err != nil {
			break
		}
		s = temp + s
	}
	return s
}

func (pc *Preprocessor) IfCondition() bool {
	if !pc.expect(IF) && !pc.expect(ELSEIF){ return false }
	if !pc.expect(LPAREN) { return false }
	if !pc.Condition() { return false }
	if !pc.expect(RPAREN) { return false }

	//get new condition
	pc.ts.PushN(pc.temp.buf)
	pc.temp.Clear()

	//deal "{"  "}"
	if pc.tok != LCB { return false }
	//replace "{" with "then"
	pc.ts.Pop()
	pc.ts.Push(strThen)
	pc.temp.Clear()

	if !pc.ScanIfBody() { return false }
	tok, lit := pc.scanIgnoreWhitespace()
	if tok == ELSEIF {
		pc.ts.Push(lit)
		pc.tok, pc.lit = tok, lit
		if !pc.IfCondition() {
			return false
		}
	}else if tok == ELSE {
		pc.ts.Push(lit)
		if !pc.ScanElseBody() { return false }
	} else if tok == EOF {
		pc.ts.Push(strEnd)
		return true
	}else { //if后面没有 else 也没有elseif
		pc.ts.Push(strEnd)
		pc.ts.Push(lit)
		if tok == IF {
			pc.tok, pc.lit = tok, lit
			if !pc.IfCondition() {
				return false
			}
		} else if tok == RCB { //end of elseif  or end of else
			pc.ts.Pop()
			lit := lit[ : len(lit)-1]
			//keep the format and add strEnd
			pc.ts.Push(lit)
			pc.ts.Push(strEnd)
		}
	}
	return true
}

func (pc *Preprocessor) ScanIfBody() bool{
	for {
		tok, lit := pc.s.Scan()
		if tok == EOF {
			return false
		}
		if tok == RCB {
			break
		}
		pc.ts.Push(lit)
		if tok == IF || tok == ELSEIF{
			pc.tok, pc.lit = tok, lit
			if !pc.IfCondition() {
				return false
			}
			if pc.tok == RCB { //when elseif condition contains if condition,it will run here
				pc.ts.Pop()
				break
			}
		}
	}
	return true
}

func (pc *Preprocessor) ScanElseBody() bool {
	tok, lit := pc.scanIgnoreWhitespace()
	if tok != LCB {
		return false
	}
	for {
		tok, lit = pc.s.Scan()
		if tok == RCB {
			break
		}
		if tok == EOF {
			return false
		}
		pc.ts.Push(lit)
		if tok == IF {
			pc.tok, pc.lit = tok, lit
			if !pc.IfCondition() {
				return false
			}
		}
	}
	pc.ts.Push(strEnd)
	return true
}
func (pc *Preprocessor) Condition() bool {
	if pc.tok == EXIST {
		if !pc.Exist() {
			return false
		}
	} else {
		if !pc.Comparison() {
			return false
		}
	}

	return true
}

func (pc *Preprocessor) Exist() bool {
	if !pc.expect(EXIST) {
		return false
	}
	if !pc.expect(LPAREN) {
		return false
	}
	if !pc.TableValue() {
		return false
	}
	if !pc.expect(RPAREN) {
		return false
	}

	return true
}

func (pc *Preprocessor) Comparison() bool {
	if !pc.SingleValue("left") {
		return false
	}
	if pc.accept(DEQUAL) {
		if !pc.SingleValue("right") {
			return false
		}
	} else if pc.accept(IN) {
		if !pc.InCondition() {
			return false
		}
		if !pc.ListLiteral() {
			return false
		}
	} else {
		pc.Err = fmt.Errorf("found %q, expected \"==\" or \"in\"", pc.lit)
		return false
	}
	return true
}

func (pc *Preprocessor) InCondition() bool {
	pc.temp.PopN(2) //pop IN
	field, err := pc.temp.PopN(pc.temp.length - 1)
	if err != nil {
		return false
	}
	pc.temp.Push(FunFieldIn)
	pc.temp.PushN(field)
	pc.temp.Push(strComma)
	return true
}

func (pc *Preprocessor) ListLiteral() bool {
	if !pc.expect(LPAREN) {
		return false
	}
	if !pc.expect(QUOTEDLIT) {
		return false
	} // first element of list
	for {
		if pc.accept(RPAREN) {
			break
		}

		if !pc.expect(COMMA) {
			return false
		}
		if !pc.expect(QUOTEDLIT) {
			return false
		}
	}
	return true
}

func (pc *Preprocessor) SingleValue(l_or_r string) bool {
	switch pc.tok {
	case QUOTEDLIT:
		pc.accept(QUOTEDLIT)
	case THIS:
		if !pc.ThisExpr() {
			return false
		}
	default:
		pc.Err = fmt.Errorf("found %q, expected double quote or \"this\"", pc.lit)
		return false
	}
	return true
}

func (pc *Preprocessor) ThisExpr() bool {
	if !pc.expect(THIS) {
		return false
	}
	if !pc.expect(DOT) {
		return false
	}

	if !pc.Field() {
		return false
	}

	if pc.accept(DOT) {
		if !pc.ParentField() {
			return false
		}
	}
	return true
}

func (pc *Preprocessor) expect(tok Token) bool {
	if pc.accept(tok) {
		return true
	}
	pc.Err = fmt.Errorf("found \"%s\", expected \"%s\"", pc.lit, tokenDisplay[int(tok)])
	return false
}

func (pc *Preprocessor) accept(tok Token) bool {

	if pc.tok == tok {
		pc.nextSym()
		pc.temp.Push(pc.lit)
		return true
	}
	return false
}

func (pc *Preprocessor) nextSym() {
	tok, lit := pc.s.Scan()
	if tok == WS {
		tok, lit = pc.s.Scan()
	}
	pc.tok, pc.lit = tok, lit
}

func (pc *Preprocessor) Field() bool {

	if !pc.accept(IDENT) {
		return false
	}
	return true
}

func (pc *Preprocessor) ParentField() bool {
	if !pc.expect(PARENT) {
		return false
	}
	if !pc.expect(DOT) {
		return false
	}
	if !pc.Field() {
		return false
	}
	if pc.accept(DOT) {
		if !pc.ParentField() {
			return false
		}
	}
	return true
}

func (pc *Preprocessor) TableValue() bool {
	pc.temp.Pop() //key "table" is useless
	if !pc.expect(TABLE) {
		return false
	}
	pc.temp.Pop()
	if !pc.expect(DOT) {
		return false
	}
	if !pc.TableName() {
		return false
	}
	for {
		if pc.tok != DOT {
			break
		}
		pc.temp.Pop()
		pc.temp.Push(strComma)
		pc.expect(DOT)
		if !pc.Where() {
			return false
		}
	}
	return true
}

func (pc *Preprocessor) TableName() bool {
	if pc.tok == IDENT {
		pc.temp.Pop()
		pc.temp.Push(`"` + pc.lit + `"`)
	}
	if pc.accept(IDENT) {
		return true
	}
	return false
}

func (pc *Preprocessor) Where() bool {
	pc.temp.Pop()
	if !pc.expect(WHERE) {
		return false
	}
	pc.temp.Pop()
	if !pc.expect(LPAREN) {
		return false
	}
	if pc.tok != IDENT {
		return false
	} else {
		pc.temp.Pop()
		pc.temp.Push(`"` + pc.lit + `"`)
	}

	if !pc.Field() {
		return false
	}

	if pc.tok == DEQUAL {
		pc.temp.Pop()
		pc.temp.Push(strDEQUAL)
	}
	if !pc.expect(DEQUAL) {
		return false
	}
	if !pc.SingleValue("right") {
		return false
	}
	pc.temp.Pop()
	if !pc.expect(RPAREN) {
		return false
	}

	return true
}

func (pc *Preprocessor) scanIgnoreWhitespace() (Token, string) {
	//Reserve ws to keep format
	ws := ""
	tok, lit := pc.s.Scan()
	if tok == WS {
		ws = lit
		tok, lit = pc.s.Scan()
	}
	pc.tok, pc.lit = tok, lit
	return tok, ws + lit
}
