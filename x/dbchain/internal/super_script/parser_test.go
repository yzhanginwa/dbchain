package super_script 

import (
    "errors"
    "reflect"
    "strings"
    "testing"
    "github.com/yzhanginwa/dbchain-sm/x/dbchain/internal/utils"
    "github.com/yzhanginwa/dbchain-sm/x/dbchain/internal/super_script/eval"
)

func TestParser_ParseConditioon(t *testing.T) {
    var tests = []struct {
        s    string
        err  string
    }{
        {
            s: `this.corp_id.parent.created_id == this.created_id`,
            err: "",
        },

        {
            s: `this.name in ("foo", "bar")`,
            err: "",
        },

        {
            s: `exist(table.corp.where(id == this.corp_id))`,
            err: "",
        },

        {
            s: `exist(table.corp.where(type == "corp").where(name == "Microsoft"))`,
            err: "",
        },

        // Errors
        {s: `this.corp_id.aa`, err: `found "aa", expected "parent"`},
        {s: `foo`, err: `found "foo", expected double quote or "this"`},
        {s: `this !`, err: `found "!", expected "dot"`},
        {s: `this field`, err: `found "field", expected "dot"`},
    }

    for i, tt := range tests {
        parser := NewParser(strings.NewReader(tt.s),
            func(table, field string) bool {
                return true
            },
            func(table, field string) (string, error) {
                return "foo", nil
            },
        )
        parser.prepareParsing()
        parser.Condition(&(eval.IfCondition{}))
        err := parser.err
        if !reflect.DeepEqual(tt.err, errstring(err)) {
            t.Errorf("%d. %q: error mismatch:\n  exp=%s\n  got=%s\n\n", i, tt.s, tt.err, err)
        }
    }
}

func TestParser_ParseExistCondition(t *testing.T) {
    var tests = []struct {
        s    string
        err  string
    }{
        {
            s: `if (exist(table.corp.where(name == "aa"))) { return(true) }`,
            err: "",
        },
        {
            s: `if (exist(table.corp.where(name == "aa"))) {
                    return(true)
                } else {
                  if(exist(table.foo.where(name == "bb"))) {
                      return(false)
                  }
                }`,
            err: "",
        },

    }

    for i, tt := range tests {
        parser := NewParser(strings.NewReader(tt.s),
            func(table, field string) bool {
                return true
            },
            func(table, field string) (string, error) {
                return "foo", nil
            },
        )
        parser.ParseFilter()
        err := parser.err
        if !reflect.DeepEqual(tt.err, errstring(err)) {
            t.Errorf("%d. %q: error mismatch:\n  exp=%s\n  got=%s\n\n", i, tt.s, tt.err, err)
        }

       if len(parser.syntaxTree) != 1 {
           t.Errorf("syntax tree error")
       }
       if len(parser.syntaxTree[0].IfCondition.Condition.Exist.TableValue.Items) != 2 {
           t.Errorf("syntax tree error")
       }


       tableName := parser.syntaxTree[0].IfCondition.Condition.Exist.TableValue.Items[0].(string)
       theWhere  := parser.syntaxTree[0].IfCondition.Condition.Exist.TableValue.Items[1].(eval.Where)
       if tableName != "corp" {
           t.Errorf("syntax tree error")
       }
       if theWhere.Field != "name" {
           t.Errorf("syntax tree error")
       }
       if theWhere.Operator != "==" {
           t.Errorf("syntax tree error")
       }
       if theWhere.Right.QuotedLit != "aa" {
           t.Errorf("syntax tree error")
       }
    }
}

func TestParser_ParseScript(t *testing.T) {
    script := `if(this.corp_id.parent.created_by == this.created_by) {
                insert("corp", "name", "foo", "mailing", "100 main st")
                insert("corp", "name", "bar", "mailing", "110 main st")
                return(false)
                } else {
                return(true)
                }
                insert("corp", "name", "bar1", "mailing", "111 main st")
                insert("corp", "name", "bar2", "mailing", this.mailing)
               `

    parser := NewParser(strings.NewReader(script),
        func(table, field string) bool {
            return true
        },
        func(table, field string) (string, error) {
            if tn, ok := utils.GetTableNameFromForeignKey(field); ok {
                return tn, nil
            } else {
                return "", errors.New("Wrong reference field name")
            }
        },
    )
    parser.ParseTrigger()
    if parser.err != nil {
        t.Errorf("Failed to parse script")
    }
    if len(parser.syntaxTree) != 3 {
        t.Errorf("syntax tree error")
    }
    if len(parser.syntaxTree[0].IfCondition.IfStatements) != 3 {
        t.Errorf("syntax tree error")
    }
    if parser.syntaxTree[0].IfCondition.Condition.Comparison.Left.ThisExpr.Items[0] != "corp_id" {
        t.Errorf("syntax tree error")
    }
    pt := parser.syntaxTree[0].IfCondition.Condition.Comparison.Left.ThisExpr.Items[1].(eval.ParentField)
    if pt.ParentTable != "corp" {
        t.Errorf("syntax tree error")
    }
    if pt.Field!= "created_by" {
        t.Errorf("syntax tree error")
    }
    if parser.syntaxTree[0].IfCondition.IfStatements[2].Return != "false" {
        t.Errorf("syntax tree error")
    }
    if len(parser.syntaxTree[0].IfCondition.ElseStatements) != 1 {
        t.Errorf("syntax tree error")
    }
    if parser.syntaxTree[0].IfCondition.ElseStatements[0].Return != "true" {
        t.Errorf("syntax tree error")
    }
    if parser.syntaxTree[2].Insert.TableName != "corp" {
        t.Errorf("syntax tree error")
    }
    if parser.syntaxTree[2].Insert.Value["name"].QuotedLit != "bar2" {
        t.Errorf("syntax tree error")
    }
    if parser.syntaxTree[2].Insert.Value["mailing"].ThisExpr.Items[0] != "mailing" {
        t.Errorf("syntax tree error")
    }

}

// errstring returns the string representation of an error.
func errstring(err error) string {
    if err != nil {
        return err.Error()
    }
    return ""
}
