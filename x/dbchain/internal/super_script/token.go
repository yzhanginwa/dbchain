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
    DQUOTE   // "
    COMMA    // ,
    DOT      // .
    LPAREN   // (
    RPAREN   // )

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
    "identity", "double quote", "comma",
    "dot" , "left parenthesis", "right parenthesis",
    "this", "parent", "table",
    "=" , "==", "in",
    "where",
}
