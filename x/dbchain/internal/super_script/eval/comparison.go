package eval

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
