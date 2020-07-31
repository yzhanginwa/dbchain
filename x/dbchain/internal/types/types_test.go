package types

import (
    "testing"
    "github.com/stretchr/testify/require"
)

func TestMetaNames(t *testing.T) {
    cases := []struct {
        valid bool
        name string
    }{
        { true, "abcd1234" },
        { false, "a-b_c-1" },
        { false, "Abcd1234" },
        { false, "aBcd1234" },
        { false, "1abce" },
        { false, "a cd" },
        { false, "a:cd" },
    }

    for _, tc := range cases {
        result := validateMetaName(tc.name)
        require.Equal(t, result, tc.valid)
    }
}

func TestEnumOption(t *testing.T) {
    cases := []struct {
        valid bool
        enumOption string
    }{
        { true, `enum("aaa")` },
        { true, `enum("aaa","bbb")` },
        { true, `enum("aaa", "b-b_b" ,"c cc")` },
        { false, `enum("aaa)` },
        { false, `enum("aaa"` },
    }

    for _, tc := range cases {
        result := ValidateEnumColumnOption(tc.enumOption)
        require.Equal(t, result, tc.valid)
    }
}
