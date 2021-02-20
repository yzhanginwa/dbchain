package super_script

import (
    "bufio"
    "bytes"
    "io"
    "strings"
)

// eof represents a marker rune for the end of the reader.
var eof = rune(0)

// Scanner represents a lexical scanner.
type Scanner struct {
    r *bufio.Reader
}

func NewScanner(r io.Reader) *Scanner {
    return &Scanner{r: bufio.NewReader(r)}
}

func (s *Scanner) Scan() (tok Token, lit string) {
    ch := s.read()

    if isWhitespace(ch) {
        s.unread()
        return s.scanWhitespace()
    } else if ch == '=' {
        secondEqual := s.read() 
        if secondEqual == '=' {
            return DEQUAL, "=="
        } else {
            s.unread() 
            return EQUAL, "="
        }
    } else if ch == '"' {
        s.unread()
        return s.scanQuotedLit()
    } else if isValidLetter(ch) {
        s.unread()
        return s.scanIdent()
    } else if isDigit(ch) {
        s.unread()
        return s.scanDigit()
    } else if ch == '!' {
        secondEqual := s.read()
        if secondEqual == '=' {
            return UNEQUAL, " ~= "
        } else {
            s.unread()
            return ILLEGAL, "!"
        }
        return
    }

    switch ch {
    case eof:
        return EOF, ""
    case ',':
        return COMMA, string(ch)
    case '.':
        return DOT, string(ch)
    case '(':
        return LPAREN, string(ch)
    case ')':
        return RPAREN, string(ch)
    case '{':
        return LCB, string(ch)
    case '}':
        return RCB, string(ch)
    }

    return ILLEGAL, string(ch)
}

func (s *Scanner) scanWhitespace() (tok Token, lit string) {
    var buf bytes.Buffer
    buf.WriteRune(s.read())

    for {
        if ch := s.read(); ch == eof {
            break
        } else if !isWhitespace(ch) {
            s.unread()
            break
        } else {
            buf.WriteRune(ch)
        }
    }

    return WS, buf.String()
}

func (s *Scanner) scanQuotedLit() (tok Token, lit string) {
    var buf bytes.Buffer
    ch := s.read()
    if ch != '"' {
        return ILLEGAL, string(ch)  // this should not happen
    }
    buf.WriteRune(ch)

    withEscape := false
    for {
        ch = s.read()
        if ch == eof {
            return ILLEGAL, string(ch)
        }
        if withEscape {
            withEscape = false
        } else {
            if ch == '\\' {
                withEscape = true
            } else if ch == '"' {
                buf.WriteRune(ch)
                break
            }
        }
        buf.WriteRune(ch)
    }

    return QUOTEDLIT, buf.String()
}

func (s *Scanner) scanIdent() (tok Token, lit string) {
    var buf bytes.Buffer
    buf.WriteRune(s.read())

    for {
        if ch := s.read(); ch == eof {
            break
        } else if !isValidLetter(ch) && !isDigit(ch) && ch != '_' {
            s.unread()
            break
        } else {
            _, _ = buf.WriteRune(ch)
        }
    }

    switch strings.ToUpper(buf.String()) {
    case "THIS":
        return THIS, buf.String()
    case "PARENT":
        return PARENT, buf.String()
    case "TABLE":
        return TABLE, buf.String()
    case "==":
        return DEQUAL, buf.String()
    case "IN":
        return IN, buf.String()
    case "IF":
        return IF, buf.String()
    case "ELSE":
        return ELSE, buf.String()
    case "INSERT":
        return INSERT, buf.String()
    case "RETURN":
        return RETURN, buf.String()
    case "TRUE":
        return TRUE, buf.String()
    case "FALSE":
        return FALSE, buf.String()
    case "WHERE":
        return WHERE, buf.String()
    case "EXIST":
        return EXIST, buf.String()
    case "ELSEIF":
        return ELSEIF, buf.String()
    case "!=":
        return UNEQUAL, buf.String()
    case "FUNCTION":
        return FUNCTION, buf.String()


    }

    // Otherwise return as a regular identifier.
    return IDENT, buf.String()
}

func (s *Scanner) scanDigit() (tok Token, lit string){
    var buf bytes.Buffer
    buf.WriteRune(s.read())

    for {
        if ch := s.read(); ch == eof {
            break
        } else if !isDigit(ch){
            s.unread()
            break
        } else {
            _, _ = buf.WriteRune(ch)
        }
    }
    return NUMBER, buf.String()
}

func (s *Scanner) read() rune {
    ch, _, err := s.r.ReadRune()
    if err != nil {
        return eof
    }
    return ch
}

func (s *Scanner) unread() { _ = s.r.UnreadRune() }

func isWhitespace(ch rune) bool {
    return ch == ' ' || ch == '\t' || ch == '\n'
}

func isValidLetter(ch rune) bool {
    return (ch >= 'a' && ch <= 'z') ||
           (ch >= 'A' && ch <= 'Z') ||
           (ch == '_')
}

func isDigit(ch rune) bool { return (ch >= '0' && ch <= '9') }

