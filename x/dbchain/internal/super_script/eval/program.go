package eval

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
