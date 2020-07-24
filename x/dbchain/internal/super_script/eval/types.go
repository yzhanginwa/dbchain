package eval

import (
    "strconv"
)

type ReturnValue int
const (
    NIL ReturnValue = iota
    FALSE
    TRUE
)

type Program struct {
    CurrentTable string
    NewRecord    map[string]string
    Script       string
    SyntaxTree   []Statement
    Return       ReturnValue
    InsertFunc   func(string, map[string]string)
    GetFieldValueFunc func(string, uint, string) string // appId, tableName, id, fieldName
}

func NewProgram(tableName string, newRecord map[string]string, script string) *Program {
    return &Program{
        CurrentTable: tableName,
        NewRecord:    newRecord,
        Script:       script,
        Return:       NIL,
    }
}

func (p *Program) EvaluateScript(syntaxTree []Statement) bool {
    for _, statement := range syntaxTree {
        statement.Evaluate(p)
        if p.Return == FALSE {
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

    if s.Insert.TableName != "" {
        s.Insert.Evaluate(p)
    }

    if len(s.IfCondition.Statements) != 0 {
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
            if p.Return == FALSE {
                return
            } else if p.Return == TRUE {
                return
            }
        }
    }
}

type Condition struct {
    Type       string
    Exist      Exist
    Comparison Comparison
}

func (c *Condition) Evaluate(p *Program) bool {
    switch c.Type {
    case "exist":
        return c.Exist.Evaluate(p)
    case "comparison":
        return c.Comparison.Evaluate(p)
    default:
        return false
    }
}

type Exist struct {
    TableValue TableValue
}

func (e *Exist) Evaluate(p *Program) bool {
    if len(e.TableValue.Evaluate(p)) > 0 {
        return true
    }
    return false
}

type Comparison struct {
    Left SingleValue
    Operator string                // "==" or "in"
    Right interface{}              // single or list value
}

func (c *Comparison) Evaluate(p *Program) bool {
    left  := c.Left.Evaluate(p)
    if c.Operator == "==" {
        right := c.Right.(SingleValue)
        rightValue := right.Evaluate(p)
        if left == rightValue {
            return true
        }
    } else if c.Operator == "in" {
        right := c.Right.(MultiValue)
        rightValue := right.Evaluate(p)
        for _, v := range rightValue {
            if left == v {
                return true
            }
        }
    }
    return false
}

type Insert struct {
    TableName string
    Value map[string]string
}

func (i *Insert) Evaluate(p *Program) {
    p.InsertFunc(i.TableName, i.Value)
}

type SingleValue struct {
    QuotedLit string
    ThisExpr ThisExpression
}

func (s *SingleValue) Evaluate(p *Program) string {
    if s.QuotedLit != "" {
        return s.QuotedLit
    }
    return s.ThisExpr.Evaluate(p)
}

type ThisExpression struct {
    Items []interface{}
}

func (t *ThisExpression) Evaluate(p *Program) string {
    currentTable := ""
    currentField := ""
    currentValue := ""
    for _, item := range t.Items {
        if currentField == "" {
            currentField = item.(string)
            currentValue = p.NewRecord[currentField]
        } else {
            parentField := item.(ParentField)
            currentTable = parentField.ParentTable
            currentField = parentField.Field
            id, _ := strconv.Atoi(currentValue)
            currentValue = p.GetFieldValueFunc(currentTable, uint(id), currentField)
        }

    }
    return currentValue
}

type ParentField struct {
    ParentTable string
    Field string
}

type MultiValue struct {
    TableValue TableValue
    ListLiteral ListLiteral
}

func (m *MultiValue) Evaluate(p *Program) []string {
    return []string{}
}

type TableValue struct {
    Items []interface{}
}

func (t *TableValue) Evaluate(p *Program) []string {
    return []string{}
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
