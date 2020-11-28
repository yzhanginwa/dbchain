package eval

type Insert struct {
    TableName string
    Value map[string]SingleValue
}

func (i *Insert) Evaluate(p *Program) {
    record := map[string]string{}
    for k, v := range i.Value {
        record[k] = v.Evaluate(p)
    }
    p.InsertFunc(i.TableName, record)
}
