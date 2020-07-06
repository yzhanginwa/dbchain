package super_script 

// Token represents a lexical token.
type Token int

const (
    // Special tokens
    ILLEGAL Token = iota
    EOF
    WS

    // Literals
    IDENT // main

    // Misc characters
    COMMA    // ,
    DOT      // .
    LPAREN   // (
    RPAREN   // )

    QUOTEDLIT // "xxxxxxx"

    // Keywords
    THIS
    PARENT
    TABLE
    EQUAL   // =
    DEQUAL  // ==
    IN      // in
    WHERE   // where
)

var tokenDisplay = []string{
    "illegal", "eof", "whitespace",
    "identity", "comma", "dot",
    "left parenthesis", "right parenthesis", "quoted string",
    "this", "parent", "table",
    "=" , "==", "in",
    "where",
}
