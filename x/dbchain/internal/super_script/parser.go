package super_script

import (
    "fmt"
    "io"
    "github.com/yzhanginwa/dbchain/x/dbchain/internal/super_script/eval"
)

// validate field name in current or specified table
type validateTableField func(string, string) bool
type getParentTable func(string, string) (string, error)

// Parser represents a parser.
type Parser struct {
    s   *Scanner
    tok Token  // last read token
    lit string // last read literal
    err error

    vtf validateTableField
    gpt getParentTable
    currentTable string
    currentField string

    syntaxTree interface{} // eval.Condition or eval.
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader, vtf validateTableField, gpt getParentTable) *Parser {
    return &Parser{s: NewScanner(r), vtf: vtf, gpt: gpt}
}

func (p *Parser) Start() {
    p.nextSym()
}

func (p *Parser) Trigger() error {
    p.syntaxTree = []eval.Block{}
    for {
        if p.tok == EOF {
            return nil
        }
        if !p.Statement() {
            break
        }
    }
    if p.err != nil {
        return p.err
    }
    return nil
}

func (p *Parser) Statement() bool {
    switch p.tok {
    case IF:
        if !p.IfCondition() { return false }
    case RETURN:
        if !p.Return() { return false }
    case INSERT:
        if !p.Insert() { return false }
    default:
        p.err = fmt.Errorf("found %q, expected if condition or insert statement", p.lit)
        return false
    }
    return true
}

func (p *Parser) Return() bool {
    // return(true)
    if !p.expect(RETURN) { return false }
    if !p.expect(LPAREN) { return false }
    if !p.accept(TRUE) {
        if !p.accept(FALSE) {
            p.err = fmt.Errorf("found %q, expected true or false", p.lit)
            return false
        }
    }
    if !p.expect(RPAREN) { return false }
    return true
}

func (p *Parser) Insert() bool {
    // insert(tableName, field1, "value1", fields2, "value2)
    if !p.expect(INSERT) { return false }
    if !p.expect(LPAREN) { return false }
    if !p.expect(QUOTEDLIT) { return false } // tableName
    fieldValuePairs := 0
    for {
        if p.accept(RPAREN) {
            break
        }
        if !p.expect(COMMA) { return false }
        if !p.expect(QUOTEDLIT) { return false }
        if !p.expect(COMMA) { return false }
        if !p.expect(QUOTEDLIT) { return false }
        fieldValuePairs += 1
    }
    return true
}

func (p *Parser) IfCondition() bool {
    if !p.expect(IF) { return false }
    p.FilterCondition()
    if !p.expect(THEN) { return false }
    p.Statement()
    for {
        if p.accept(FI) {
            break
        }
        if !p.Statement() {
            return false
        }
    }
    return true
}

func (p *Parser) FilterCondition() bool {
    if ! p.SingleValue() { return false }

    if p.accept(DEQUAL) {
        if !p.SingleValue() { return false }
    } else if p.accept(IN) {
        if !p.MultiValue() { return false }
    } else {
        p.err = fmt.Errorf("found %q, expected \"==\" or \"in\"", p.lit)
        return false
    }
    return true
}

func (p *Parser) SingleValue() bool {
    switch p.tok {
    case QUOTEDLIT:
        return true
    case THIS:
        if !p.ThisExpr() { return false }
    default:
        p.err = fmt.Errorf("found %q, expected double quote or \"this\"", p.lit)
        return false
    }

    return true
}

func (p *Parser) ThisExpr() bool {
    if !p.expect(THIS) { return false }
    p.currentTable = ""
    if !p.expect(DOT) { return false }
    if !p.Field() { return false }
    if p.accept(DOT) {
        if !p.ParentField() { return false }
    }
    return true
}

func (p *Parser) MultiValue() bool {
    switch p.tok {
    case TABLE:
        p.TableValue()
    case LPAREN:
        p.ListValue()
    default:
        p.err = fmt.Errorf("table or list is wanted")
        return false
    }

    return true
}

func (p *Parser) ListValue() bool {
    if !p.expect(LPAREN) { return false }
    if !p.expect(QUOTEDLIT) { return false } // first element of list
    for {
        if p.accept(RPAREN) {
            break;
        }

        if !p.expect(COMMA) { return false }
        if !p.expect(QUOTEDLIT) { return false }
    }
    return true
}

func (p *Parser) TableValue() bool {
    if !p.expect(TABLE) { return false }
    if !p.expect(DOT) { return false }
    if !p.TableName() { return false }
    if !p.expect(DOT) { return false }
    for ok := p.accept(WHERE); ok;  ok = p.accept(WHERE) {
        if !p.expect(DOT) { break }
    } 
    if !p.Field() { return false }
    return true
}

func (p *Parser) ParentField() bool {
    if !p.expect(PARENT) { return false }
    tn, err := p.gpt(p.currentTable, p.currentField)
    if err != nil {
        p.err = err
        return false
    }
    p.currentTable = tn
    if !p.expect(DOT) { return false }
    if !p.Field() { return false }
    if p.accept(DOT) {
        if !p.ParentField() { return false }
    }
    return true
}

func (p *Parser) TableName() bool {
    if p.accept(IDENT) {
        p.currentTable = p.lit
        return true
    }
    return false
}

func (p *Parser) Where() bool {
    if !p.expect(WHERE) { return false }
    if !p.expect(LPAREN) { return false }
    if !p.Field() { return false }
    if !p.expect(DEQUAL) { return false }
    if !p.SingleValue() { return false }
    if !p.expect(LPAREN) { return false }
    return true
}

func (p *Parser) Field() bool {
    if p.accept(IDENT) {
        if p.vtf(p.currentTable, p.currentField) {
            return true
        } else {
            p.err = fmt.Errorf("Field name does not exist")
            return false
        }
    }
    return false
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) nextSym(){
    tok, lit := p.scanIgnoreWhitespace()
    p.tok, p.lit = tok, lit
}

func (p *Parser) accept(tok Token) bool {
    if (p.tok == tok) {
        p.nextSym()
        return true
    }
    return false;
}

func (p *Parser) expect(tok Token) bool {
    if (p.accept(tok)) {
        return true;
    }
    p.err = fmt.Errorf("found \"%s\", expected \"%s\"", p.lit, tokenDisplay[int(tok)])
    return false;
}


// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (Token, string) {
    tok, lit := p.s.Scan()
    if tok == WS {
        tok, lit = p.s.Scan()
    }
    return tok, lit
}
