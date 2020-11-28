package eval

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
