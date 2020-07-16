package super_script 

import (
    "strings"
    "testing"
)

// Ensure the scanner can scan tokens correctly.
func TestScanner_Scan(t *testing.T) {
    var tests = []struct {
        s   string
        tok Token
        lit string
    }{
        // Special tokens (EOF, ILLEGAL, WS)
        {s: ``, tok: EOF},
        {s: `#`, tok: ILLEGAL, lit: `#`},
        {s: ` `, tok: WS, lit: " "},
        {s: "\t", tok: WS, lit: "\t"},
        {s: "\n", tok: WS, lit: "\n"},

        // Misc characters
        {s: "(", tok: LPAREN, lit: "("},
        {s: ")", tok: RPAREN, lit: ")"},

        // Identifiers
        {s: `foo`, tok: IDENT, lit: `foo`},
        {s: `Zx12_3U_-`, tok: IDENT, lit: `Zx12_3U_`},

        {s: `"abcd"`, tok: QUOTEDLIT, lit: "\"abcd\""},
        {s: `"abcd \""`, tok: QUOTEDLIT, lit: `"abcd \""`},


        // Keywords
        {s: `THIS`, tok: THIS, lit: "THIS"},
        {s: `PARENT`, tok: PARENT, lit: "PARENT"},
        {s: `=`, tok: EQUAL, lit: "="},
        {s: `==`, tok: DEQUAL, lit: "=="},
        {s: `if`, tok: IF, lit: "if"},
        {s: `then`, tok: THEN, lit: "then"},
        {s: `fi`, tok: FI, lit: "fi"},
        {s: `insert`, tok: INSERT, lit: "insert"},
        {s: `return`, tok: RETURN, lit: "return"},
        {s: `true`, tok: TRUE, lit: "true"},
        {s: `false`, tok: FALSE, lit: "false"},
    }

    for i, tt := range tests {
        s := NewScanner(strings.NewReader(tt.s))
        tok, lit := s.Scan()
        if tt.tok != tok {
            t.Errorf("%d. %q token mismatch: exp=%q got=%q <%q>", i, tt.s, tt.tok, tok, lit)
        } else if tt.lit != lit {
            t.Errorf("%d. %q literal mismatch: exp=%q got=%q", i, tt.s, tt.lit, lit)
        }
    }
}
