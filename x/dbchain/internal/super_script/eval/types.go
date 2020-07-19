package eval

import (

)

type ReturnValue int
const (
    NIL ReturnValue = iota
    FALSE
    TRUE
)

type Program struct {
    CurrentAppId uint
    CurrentTable string
    NewRecord    map[string]string
    Script       string
    SyntaxTree   []Statement
    Return       ReturnValue
}

func NewProgram(appId uint, tableName string, newRecord map[string]string, script string) *Program {
    return &Program{
        CurrentAppId: appId,
        CurrentTable: tableName,
        NewRecord:    newRecord,
        Script:       script,
        Return:       NIL
    }
}

func (p *Program) EvaluateScript(syntaxTree []Statement) bool {
    for _, statement := range syntaxTree {
        statement.Evaluate(p)
        if p.Return == FALSE
            return false
        } else if p.Return == TRUE {
            return true
        }
    }
    return true
}

type Statement struct {
    IfCondition IfCondition           // when nil for just insert statements
    Insert Insert
    Return string
}

func (s *Statement) Evaluate(p *Program) {
    switch s.Return {
    case "true":
        p.Return = TRUE
        return
    case "false":
        p.Return = FALSE
        return
    }

    if s.Insert != nil {
        (s.Insert).Evaluate(p)
    }

    if s.IfCondition != nil {
        (s.IfCondition).Evaluate(p)
    }
}

type IfCondition struct {
    Condition Condition
    Statements []Statement
}

func (ic *IfCondition) Evaluate(p *Program) {
    if (ic.Condition).Evaluate(p) {
        for _, statement := range ic.Statements{
            statement.Evaluate(p)
            if p.Return == FALSE
                return false
            } else if p.Return == TRUE {
                return true
            }
        }
    }
}

type Condition struct {
    Left SingleValue
    Operator string                // "==" or "in"
    Right interface{}              // single or multi value
}

func (c *Condition) Evaluate(p *Program) bool {
    return true
}

type Insert struct {
    TableName string
    Value map[string]string
}

func (c *Insert) Evaluate(p *Program) {

}

type SingleValue struct {
    QuotedLit string
    ThisExpr ThisExpression
}

func (s *SingleValue) Evaluate(p *Program) {

}

type ThisExpression struct {
    Items []interface{}
}

func (t *ThisExpression) Evaluate(p *Program) {

}

type MultiValue struct {
    TableValue TableValue
    ListLiteral ListLiteral
}

func (m *MultiValue) Evaluate(p *Program) {

}

type TableValue struct {
    Items []interface{}
}

func (t *TableValue) Evaluate(p *Program) {

}

type ListLiteral struct {
    Items []string
}

type Where struct {          // parent is TableValue.Items
    Field string             // field name of a table
    Operator string
    Right interface{}
}

type Field struct {
    Item string
}
