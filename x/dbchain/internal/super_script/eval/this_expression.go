package eval

import (
    "strconv"
)

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
