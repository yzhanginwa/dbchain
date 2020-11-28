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

type getFieldValueCallback func(string, uint, string) string  // tableName, id, fieldName
type getTableValueCallback func([](map[string]string)) [](map[string]string)  // input: querierObjs; outptu: rows of result set with only field "id"
type insertCallback func(string, map[string]string)

type Program struct {
    CurrentTable string
    NewRecord    map[string]string
    Script       string
    SyntaxTree   []Statement
    Return       ReturnValue

    GetFieldValueFunc getFieldValueCallback
    GetTableValueFunc getTableValueCallback
    InsertFunc        insertCallback
}

func NewProgram(tableName string, newRecord map[string]string, script string,
                fieldValueFunc getFieldValueCallback,
                tableValueFunc getTableValueCallback,
                insertFunc insertCallback) *Program {
    return &Program{
        CurrentTable: tableName,
        NewRecord:    newRecord,
        Script:       script,
        Return:       NIL,
        GetFieldValueFunc: fieldValueFunc,
        GetTableValueFunc: tableValueFunc,
        InsertFunc:        insertFunc,
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
    return false                      // if no explicit true/false, we invalidate the filter
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

    if len(s.IfCondition.IfStatements) != 0 || len(s.IfCondition.ElseStatements) != 0 {
        (s.IfCondition).Evaluate(p)
    }
}

type IfCondition struct {
    Condition Condition
    IfStatements []Statement
    ElseStatements []Statement
}

func (ic *IfCondition) Evaluate(p *Program) {
    var statements []Statement
    if (ic.Condition).Evaluate(p) {
        statements = ic.IfStatements
    } else {
        statements = ic.ElseStatements
    } 
    for _, statement := range statements {
        statement.Evaluate(p)
        if p.Return == FALSE {
            return
        } else if p.Return == TRUE {
            return
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
        result := (c.Exist).Evaluate(p)
        return result
    case "comparison":
        result := (c.Comparison).Evaluate(p)
        return result
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
        right := c.Right.([]string)
        for _, v := range right {
            if left == v {
                return true
            }
        }
    }
    return false
}

type Insert struct {
    TableName string
    Value map[string]SingleValue
}

func (i *Insert) Evaluate(p *Program) {
    record := map[string]string{}
    for k, v := range i.Value {
        record[k] = v.Evaluate(p)
    }
    p.InsertFunc(i.TableName, record)
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

type TableValue struct {
    Items []interface{}
}

func (t *TableValue) Evaluate(p *Program) []map[string]string {
    qo := map[string]string{
        "method": "table",
        "table": t.Items[0].(string),
    }
    querierObjs := []map[string]string{qo}

    for _, item := range t.Items[1:] {
        theWhere := item.(Where)
        if theWhere.Operator != "==" {
            continue
        }
        qo := map[string]string{
            "method": "equal",
            "field": theWhere.Field,
            "value": theWhere.Right.Evaluate(p),
        }
        querierObjs = append(querierObjs, qo)
    }
    result := p.GetTableValueFunc(querierObjs)
    return result
}

type ListLiteral struct {
    Items []string
}

type Where struct {          // parent is TableValue.Items
    Field string             // field name of a table
    Operator string
    Right SingleValue
}
