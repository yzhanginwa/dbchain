package eval

import (

)

type ConditionOperator int

const (
    EQUAL ConditionOperator = iota
    IN
)

type filterCondition struct {
    left singleValue
    operator ConditionOperator
    right inteface{}
}

type singleValue struct {
    quotedLit string
    thisExpr thisExpression
}

type thisExpression struct {
    items []inteface{}
}

type multiValue struct {
    items []interface{}
}

