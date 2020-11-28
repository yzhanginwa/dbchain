package eval

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
