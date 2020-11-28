package eval

type ReturnValue int
const (
    NIL ReturnValue = iota
    FALSE
    TRUE
)

type getFieldValueCallback func(string, uint, string) string  // tableName, id, fieldName
type getTableValueCallback func([](map[string]string)) [](map[string]string)  // input: querierObjs; outptu: rows of result set with only field "id"
type insertCallback func(string, map[string]string)

type ParentField struct {
    ParentTable string
    Field string
}

type ListLiteral struct {
    Items []string
}

type Where struct {          // parent is TableValue.Items
    Field string             // field name of a table
    Operator string
    Right SingleValue
}
