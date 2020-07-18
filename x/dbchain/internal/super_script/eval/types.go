package eval

import (

)

type Program struct {
    CurrentAppId uint
    CurrentTable string
    NewRecord map[string]string
    Script string
    SyntaxTree []Statement
}

func NewProgram(appId uint, tableName string, newRecord map[string]string, script string) *Program {
    return &Program{
        CurrentAppId: appId,
        CurrentTable: tableName,
        NewRecord:    newRecord,
        Script:       script,
    }
}

func (p *Program) ParseScript() {

}

type Statement struct {
    IfCondition IfCondition           // when nil for just insert statements
    Insert Insert
    Return string
}

type IfCondition struct {
    Condition Condition
    Statements []Statement
}

type Condition struct {
    Left SingleValue
    Operator string                // "==" or "in"
    Right interface{}              // single or multi value
}

type Insert struct {
    TableName string
    Value map[string]string
}

type SingleValue struct {
    QuotedLit string
    ThisExpr ThisExpression
}

type ThisExpression struct {
    items []interface{}
}

type MultiValue struct {
    TableValue TableValue
    ListLiteral ListLiteral
}

type TableValue struct {
    Items []interface{}
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
