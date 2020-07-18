package super_script 

import (
    "reflect"
    "strings"
    "testing"
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
            s: `this.corp_id in table.corp.id`,
            err: "",
        },

        {
            s: `this.corp_id in table.corp.where(type == "corp").name`,
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
        parser.Start()
        parser.Condition()
        err := parser.err
        if !reflect.DeepEqual(tt.err, errstring(err)) {
            t.Errorf("%d. %q: error mismatch:\n  exp=%s\n  got=%s\n\n", i, tt.s, tt.err, err)
        }
    }
}

func TestParser_ParseScript(t *testing.T) {
    var tests = []struct {
        s    string
        err  string
    }{
        {
            s: `if this.corp_id.parent.created_id == this.created_id then
                insert("corp", "name", "foo", "mailing", "100 main st")
                insert("corp", "name", "bar", "mailing", "110 main st")
                return(false)
                fi
                insert("corp", "name", "bar1", "mailing", "111 main st")
                insert("corp", "name", "bar2", "mailing", "112 main st")
               `,
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
        parser.Start()
        parser.Script()
        if !reflect.DeepEqual(tt.err, errstring(parser.err)) {
            t.Errorf("%d. %q: error mismatch:\n  exp=%s\n  got=%s\n\n", i, tt.s, tt.err, parser.err)
        }
    }
}

// errstring returns the string representation of an error.
func errstring(err error) string {
    if err != nil {
        return err.Error()
    }
    return ""
}
