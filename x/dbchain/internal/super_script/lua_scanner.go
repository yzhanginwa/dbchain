package super_script

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type LuaScanner struct {
	Scanner
}

func NewLuaScanner (r io.Reader) *LuaScanner {
	l := LuaScanner{}
	l.r = bufio.NewReader(r)
	return &l
}

func (ls *LuaScanner) Scan() (tok Token, lit string) {
	ch := ls.read()

	if isWhitespace(ch) {
		ls.unread()
		return ls.scanWhitespace()
	} else if ch == '"' {
		ls.unread()
		return ls.scanQuotedLit()
	} else if isValidLetter(ch) {
		ls.unread()
		return ls.scanIdent()
	} else if isDigit(ch) {
		ls.unread()
		return ls.scanDigit()
	}

	if ch == eof {
		return EOF, ""
	} else if ch ==  '('{
		return LPAREN, string(ch)
	} else if ch == ')' {
		return RPAREN, string(ch)
	} else {
		return OTHER, string(ch)
	}
}

func (ls *LuaScanner) scanIdent() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(ls.read())

	for {
		if ch := ls.read(); ch == eof {
			break
		} else if !isValidLetter(ch) && !isDigit(ch) && ch != '_' {
			ls.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	switch strings.ToUpper(buf.String()) {

	case "FUNCTION":
		return FUNCTION, buf.String()
	case "FOR":
		return FOR, buf.String()
	case "WHILE":
		return WHILE, buf.String()
	case "REPEAT":
		return REPEAT, buf.String()
	case "DO":
		return DO, buf.String()
	case "__SCRIPT_LOOP_COUNT__":
		return SCRIPT_LOOP_COUNT, buf.String()
	}

	// Otherwise return as a regular identifier.
	return IDENT, buf.String()

}