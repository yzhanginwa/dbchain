package eval

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
