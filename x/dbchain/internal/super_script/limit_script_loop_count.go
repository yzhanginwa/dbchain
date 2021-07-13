package super_script

import (
	"io"
)

type EmbedLoopCount struct {
	ts      *tokenStack
	s       *LuaScanner
	tok     Token  // last read token
	lit     string // last read literal
	Success bool
	hasEmbedLoopCountSymbol bool
}

const consumeGas = `
	scriptConsumeGas()
`

func NewPreprocessor(r io.Reader) *EmbedLoopCount {
	return &EmbedLoopCount{
		s:    NewLuaScanner(r),
		ts:   newTokenStack(),
		Success : true,
	}
}

func (elc *EmbedLoopCount) Process(){

	//to check it is a chunk or function, there may be loop both in chunk and function
	if !elc.scanFirstTok(){
		return
	}

	for {
		tok, lit := elc.s.Scan()
		if tok == EOF {// end of scan
			break
		} else if tok == ILLEGAL {
			elc.ts.Clear()
			elc.Success =false
			return
		}
		elc.ts.Push(lit)
		if tok == FUNCTION {
			elc.tok, elc.lit = tok, lit
			if !elc.FuncHandler() {
				return
			}
		} else if tok == FOR || tok == WHILE {
			elc.tok, elc.lit = tok, lit
			if !elc.LoopHandler() {
				return
			}
		} else if tok == REPEAT {
			elc.ts.Push(consumeGas)
		}
	}
}

func (elc *EmbedLoopCount) scanFirstTok() bool {
	firstTok, lit := elc.s.Scan()
	if firstTok == WS {
		firstTok, lit = elc.s.Scan()
	}
	if firstTok == FUNCTION {
		elc.ts.Push(lit)
		elc.tok, elc.lit = firstTok, lit
		if !elc.FuncHandler() {
			return false
		}
		return true
	} else {
		elc.embedCountSymbol()
	}

	elc.ts.Push(lit)
	if firstTok == FOR || firstTok == WHILE {
		elc.tok, elc.lit = firstTok, lit
		if !elc.LoopHandler() {
			return false
		}
	}

	return true
}

func (elc *EmbedLoopCount) FuncHandler() bool {
	if !elc.FuncStart() {
		elc.ts.Clear()
		elc.Success = false
		return false
	} else if elc.tok == FOR || elc.tok == WHILE {		//this condition is used to deal function begin with loop
		if !elc.LoopCondition() {
			elc.ts.Clear()
			elc.Success = false
			return false
		}
	}
	return true
}

func (elc *EmbedLoopCount) LoopHandler() bool {
	if !elc.LoopCondition() {
		elc.ts.Clear()
		elc.Success = false
		return false
	}
	return true
}
// format function test(a,b,c)
func (elc *EmbedLoopCount) FuncStart() bool {
	if !elc.expect(FUNCTION) {
		return false
	}
	if !elc.expect(IDENT) {
		return false
	}
	if !elc.scanParentheses() {
		return false
	}

	return elc.embedCountSymbol()
}

func (elc *EmbedLoopCount) embedCountSymbol() bool {
	if elc.hasEmbedLoopCountSymbol {
		return true
	}
	previousIdent, err := elc.ts.Pop()
	if err != nil {
		//elc.ts.Push(consumeGas) it will be consume gas only in loop
		elc.hasEmbedLoopCountSymbol = true
		return true
	}
	//elc.ts.Push(consumeGas)  it will be consume gas only in loop
	elc.ts.Push(previousIdent)
	elc.hasEmbedLoopCountSymbol = true
	return true
}

func (elc *EmbedLoopCount) expect(tok Token) bool {
	if elc.accept(tok) {
		return true
	}
	return false
}

func (elc *EmbedLoopCount) accept(tok Token) bool {

	if elc.tok == tok {
		elc.nextSym()
		elc.ts.Push(elc.lit)
		return true
	}
	return false
}

func (elc *EmbedLoopCount) nextSym() {
	tok, lit := elc.s.Scan()
	strWS := ""
	if tok == WS {
		strWS = lit
		tok, lit = elc.s.Scan()
	}
	elc.tok, elc.lit = tok, strWS + lit
}

func (elc *EmbedLoopCount) scanParentheses() bool {
	if !elc.expect(LPAREN) {
		return false
	}
	for {
		if elc.tok == LPAREN {
			if !elc.scanParentheses() {
				return false
			}
		} else if elc.expect(RPAREN) {
			return true
		}
		elc.nextSym()
		if elc.tok == SCRIPT_LOOP_COUNT {
			return false
		}
		elc.ts.Push(elc.lit)
	}

	return true
}

func (elc *EmbedLoopCount) LoopCondition() bool {
	if !elc.expect(FOR) && !elc.expect(WHILE){ return false }
	for {
		tok, lit := elc.s.Scan()
		elc.ts.Push(lit)
		if tok == ILLEGAL {
			return false
		} else if tok == DO {
			elc.ts.Push(consumeGas)
			return true
		} else if tok == EOF {
			return false
		}
	}
}


func (elc *EmbedLoopCount) Reconstruct() string {
	s := ""
	for {
		temp, err := elc.ts.Pop()
		if err != nil {
			break
		}
		s = temp + s
	}
	return s
}