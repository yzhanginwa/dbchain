package eval

import (

)

type ConditionOperator int

const (
    EQUAL ConditionOperator = iota
    IN
)

type Program struct {
    CurrentAppId uint
    CurrentTable string
    NewRecord map[string]string
    Script string
    Code []Block
}

func NewProgram(appId uint, tableName string, newRecord map[string]string, script string) *Program {
    return &Program{
        CurrentAppId: appId,
        CurrentTable: tableName,
        NewRecord:    newRecord,
        Script:       script,
    }
}

func (p *Program) ParseTrigger() {

}

func (p *Program) ParseFilter() {

}

type Block struct {
    Condition Condition           // when nil for just insert statements
    Insert Insert
}

type Condition struct {
    left SingleValue
    operator ConditionOperator
    right interface{}              // single or multi value
}

type Insert struct {
    tableName string
    value map[string]string
}

type SingleValue struct {
    QuotedLit string
    ThisExpr ThisExpression
}

type ThisExpression struct {
    items []interface{}
}

type MultiValue struct {
    tv TableValue
    ll ListLiteral
}

type TableValue struct {
    items []interface{}
}

type ListLiteral struct {
    items []string
}
