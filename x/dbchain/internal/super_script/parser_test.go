package super_script 

import (
    "reflect"
    "strings"
    "testing"
)

// Ensure the parser can parse strings into Statement ASTs.
func TestParser_ParseStatement(t *testing.T) {
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
        {s: `this.corp_id.aa`, err: `found "aa", expected parent`},
        {s: `foo`, err: `found "foo", expected double quote or "this"`},
        {s: `this !`, err: `found "!", expected dot`},
        {s: `this field`, err: `found "field", expected dot`},
    }

    for i, tt := range tests {
        err := NewParser(strings.NewReader(tt.s)).Comparison()
        if !reflect.DeepEqual(tt.err, errstring(err)) {
            t.Errorf("%d. %q: error mismatch:\n  exp=%s\n  got=%s\n\n", i, tt.s, tt.err, err)
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
