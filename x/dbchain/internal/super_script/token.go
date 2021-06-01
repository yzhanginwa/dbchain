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
    UNEQUAL // !=
    IN      // in
    WHERE   // where
    IF      // if
    ELSEIF  //elseif
    ELSE    // else
    INSERT  // insert
    RETURN  // return
    TRUE    // true
    FALSE   // false
    EXIST   // exist
    NUMBER  // number 123456
    FUNCTION //function
    FOR      //for
    //forbidden key words
    WHILE    //while
    REPEAT   //repeat
    SCRIPT_LOOP_COUNT //__script_loop_count__
    DO       //do
    OTHER    //other
)

var tokenDisplay = []string{
    "illegal", "eof", "whitespace",
    "identity", "comma", "dot",
    "left parenthesis", "right parenthesis", "left brace", "right brace", "quoted string",
    "this", "parent", "table",
    "=" , "==", "!=", "in",
    "where",
    "if", "elseif", "else", "insert", "return", "true", "false", "exist", "number", "function",
    "while", "repeat", "__script_loop_count__", "do", "other",
}
