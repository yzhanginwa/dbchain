package eval

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
