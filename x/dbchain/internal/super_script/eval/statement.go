package eval

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
