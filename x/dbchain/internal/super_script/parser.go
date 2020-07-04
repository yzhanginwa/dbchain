package super_script

import (
    "fmt"
    "io"
)

// Parser represents a parser.
type Parser struct {
    s   *Scanner
    tok Token  // last read token
    lit string // last read literal
    err error
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
    return &Parser{s: NewScanner(r)}
}

func (p *Parser) Comparison() error {
    p.nextSym()
    p.SingleValue()
    if p.err != nil {
        return p.err
    }

    if p.accept(DEQUAL) {
        p.SingleValue()
        if p.err != nil {
            return p.err
        }
    } else if p.accept(IN) {
        p.MultiValue()
        if p.err != nil {
            return p.err
        }
    } else {
        return fmt.Errorf("found %q, expected \"==\" or \"in\"", p.lit)
    }
    return nil
}

func (p *Parser) SingleValue() bool {
    switch p.tok {
    case DQUOTE:
        if !p.StringLiteral() { return false }
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
    if !p.expect(DOT) { return false }
    if !p.Field() { return false }
    if p.accept(DOT) {
        if !p.ParentField() { return false }
    }
    return true
}

func (p *Parser) MultiValue() bool {
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

func (p *Parser) StringLiteral() bool {
    if !p.expect(DQUOTE) { return false }
    if !p.expect(IDENT) { return false }
    if !p.expect(DQUOTE) { return false }
    return true
}

func (p *Parser) ParentField() bool {
    if !p.expect(PARENT) { return false }
    if !p.expect(DOT) { return false }
    if !p.Field() { return false }
    if p.accept(DOT) {
        if !p.ParentField() { return false }
    }
    return true
}

func (p *Parser) TableName() bool {
    if p.accept(IDENT) {
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
        return true
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
    p.err = fmt.Errorf("found \"%s\", expected %s", p.lit, tokenDisplay[int(tok)])
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
