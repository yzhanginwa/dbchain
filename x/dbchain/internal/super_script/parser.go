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

    syntaxTree []eval.Statement
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader, vtf validateTableField, gpt getParentTable) *Parser {
    return &Parser{s: NewScanner(r), vtf: vtf, gpt: gpt}
}

func (p *Parser) Start() {
    p.nextSym()
}

func (p *Parser) Script() error {
    statements := []eval.Statement{}
    for {
        if p.tok == EOF {
            return nil
        }
        if !p.Statement(&statements) {
            break
        }
    }
    if p.err != nil {
        return p.err
    }
    p.syntaxTree = statements
    return nil
}

func (p *Parser) Statement(parent *[]eval.Statement) bool {
    thisStatement := eval.Statement{}

    switch p.tok {
    case IF:
        if !p.IfCondition(&thisStatement) { return false }
    case RETURN:
        if !p.Return(&thisStatement) { return false }
    case INSERT:
        if !p.Insert(&thisStatement) { return false }
    default:
        p.err = fmt.Errorf("found %q, expected if condition or insert statement", p.lit)
        return false
    }
    (*parent) = append((*parent), thisStatement)
    return true
}

func (p *Parser) Return(parent *eval.Statement) bool {
    if !p.expect(RETURN) { return false }
    if !p.expect(LPAREN) { return false }

    returnedValue := p.lit
    if !p.accept(TRUE) {
        if !p.accept(FALSE) {
            p.err = fmt.Errorf("found %q, expected true or false", p.lit)
            return false
        }
    }

    if !p.expect(RPAREN) { return false }
    parent.Return = returnedValue
    return true
}

func (p *Parser) Insert(parent *eval.Statement) bool {
    // insert(tableName, field1, "value1", fields2, "value2)
    insert := eval.Insert{}
    if !p.expect(INSERT) { return false }
    if !p.expect(LPAREN) { return false }

    insert.TableName = p.lit
    values := make(map[string]string)

    if !p.expect(QUOTEDLIT) { return false } // tableName
    fieldValuePairs := 0
    for {
        if p.accept(RPAREN) {
            break
        }
        if !p.expect(COMMA) { return false }
        k := p.lit
        if !p.expect(QUOTEDLIT) { return false }
        if !p.expect(COMMA) { return false }
        values[k] = p.lit
        if !p.expect(QUOTEDLIT) { return false }
        fieldValuePairs += 1
    }
    insert.Value = values
    parent.Insert = insert
    return true
}

func (p *Parser) IfCondition(parent *eval.Statement) bool {
    ifCondition := eval.IfCondition{}
    statements := []eval.Statement{}

    if !p.expect(IF) { return false }
    p.Condition(&ifCondition)

    if !p.expect(THEN) { return false }
    p.Statement(&statements)
    for {
        if p.accept(FI) {
            break
        }
        if !p.Statement(&statements) {
            return false
        }
    }
    ifCondition.Statements = statements
    parent.IfCondition = ifCondition
    return true
}

func (p *Parser) Condition(parent *eval.IfCondition) bool {
    condition := eval.Condition{}

    if ! p.SingleValue(&condition, "left") { return false }

    condition.Operator = p.lit
    if p.accept(DEQUAL) {
        if !p.SingleValue(&condition, "right") { return false }
    } else if p.accept(IN) {
        if !p.MultiValue(&condition) { return false }
    } else {
        p.err = fmt.Errorf("found %q, expected \"==\" or \"in\"", p.lit)
        return false
    }
    parent.Condition = condition
    return true
}

// SingleValue could be used in Condition and Where
func (p *Parser) SingleValue(parent interface{}, l_or_r string) bool {
    singleValue := eval.SingleValue{}

    switch p.tok {
    case QUOTEDLIT:
        singleValue.QuotedLit = p.lit
        p.accept(QUOTEDLIT)
        return true
    case THIS:
        if !p.ThisExpr(&singleValue) { return false }
    default:
        p.err = fmt.Errorf("found %q, expected double quote or \"this\"", p.lit)
        return false
    }

    switch parent.(type) {
    case *eval.Condition:
        v := parent.(*eval.Condition)
        if l_or_r == "left" {
            v.Left = singleValue
        } else {
            v.Right = singleValue
        }
    case *eval.Where:
        v := parent.(*eval.Condition)
        if l_or_r == "left" {
            v.Left = singleValue
        } else {
            v.Right = singleValue
        }
    default:
        p.err = fmt.Errorf("This is impossible")
    }

    return true
}

func (p *Parser) ThisExpr(parent *eval.SingleValue) bool {
    thisExpr := eval.ThisExpression{}

    if !p.expect(THIS) { return false }
    p.currentTable = ""
    if !p.expect(DOT) { return false }

    thisExpr.Items = append(thisExpr.Items, p.lit)
    if !p.Field() { return false }

    if p.accept(DOT) {
        if !p.ParentField() { return false }
    }
    parent.ThisExpr = thisExpr
    return true
}

func (p *Parser) MultiValue(parent *eval.Condition) bool {
    multiValue := eval.MultiValue{}

    switch p.tok {
    case TABLE:
        p.TableValue(&multiValue)
    case LPAREN:
        p.ListLiteral(&multiValue)
    default:
        p.err = fmt.Errorf("table or list is wanted")
        return false
    }
    parent.Right = multiValue
    return true
}

func (p *Parser) ListLiteral(parent *eval.MultiValue) bool {
    listLiteral := eval.ListLiteral{}
    items := []string{}

    if !p.expect(LPAREN) { return false }
    items = append(items, p.lit)
    if !p.expect(QUOTEDLIT) { return false } // first element of list
    for {
        if p.accept(RPAREN) {
            break;
        }

        if !p.expect(COMMA) { return false }
        items = append(items, p.lit)
        if !p.expect(QUOTEDLIT) { return false }
    }
    listLiteral.Items = items
    parent.ListLiteral = listLiteral
    return true
}

func (p *Parser) TableValue(parent *eval.MultiValue) bool {
    tableValue := eval.TableValue{}
    items := []interface{}{}

    if !p.expect(TABLE) { return false }
    if !p.expect(DOT) { return false }
    items = append(items, p.lit)
    if !p.TableName() { return false }
    if !p.expect(DOT) { return false }
    //for ok := p.accept(WHERE); ok;  ok = p.accept(WHERE) {
    for {
        if p.tok != WHERE { break }
        if !p.Where(&items) { return false }
        //if !p.expect(DOT) { break }
    } 
    items = append(items, p.lit)
    if !p.Field() { return false }
    tableValue.Items = items
    parent.TableValue = tableValue
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

func (p *Parser) Where(parent *[]interface{}) bool {
    theWhere := eval.Where{}
 
    if !p.expect(WHERE) { return false }
    if !p.expect(LPAREN) { return false }
    theWhere.Field = p.lit
    if !p.Field() { return false }
    theWhere.Operator = p.lit
    if !p.expect(DEQUAL) { return false }
    if !p.SingleValue(&theWhere, "right") { return false }
    if !p.expect(LPAREN) { return false }
    
    (*parent) = append((*parent), theWhere)
    return true
}

func (p *Parser) Field() bool {
    fieldName := p.lit
    if p.accept(IDENT) {
        p.currentField = fieldName
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
