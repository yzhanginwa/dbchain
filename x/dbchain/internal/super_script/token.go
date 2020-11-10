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

    LCB      // {
    RCB      // }
    QUOTEDLIT // "xxxxxxx"

    // Keywords
    THIS
    PARENT
    TABLE
    EQUAL   // =
    DEQUAL  // ==
    IN      // in
    WHERE   // where
    IF      // if
    INSERT  // insert
    RETURN  // return
    TRUE    // true
    FALSE   // false
    EXIST   // exist
)

var tokenDisplay = []string{
    "illegal", "eof", "whitespace",
    "identity", "comma", "dot",
    "left parenthesis", "right parenthesis", "left brace", "right brace", "quoted string",
    "this", "parent", "table",
    "=" , "==", "in",
    "where",
}
