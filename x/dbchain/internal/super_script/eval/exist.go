package eval

type Exist struct {
    TableValue TableValue
}

func (e *Exist) Evaluate(p *Program) bool {
    if len(e.TableValue.Evaluate(p)) > 0 {
        return true
    }
    return false
}
