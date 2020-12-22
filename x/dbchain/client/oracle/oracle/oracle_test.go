package oracle

import (
    "testing"
    "github.com/stretchr/testify/require"
)

func TestMakeBatchs(t *testing.T) {
    msgs := []UniversalMsg{}
    result := makeBatches(msgs, 3)
    require.Equal(t, len(result), 0 )

    msgs = []UniversalMsg{"aaa"}
    result = makeBatches(msgs, 3)
    require.Equal(t, len(result), 1 )
    require.Equal(t, len(result[0]), 1 )
    require.Equal(t, result[0][0], "aaa" )

    msgs = []UniversalMsg{"aaa", "bbb"}
    result = makeBatches(msgs, 3)
    require.Equal(t, len(result), 1 )
    require.Equal(t, len(result[0]), 2 )
    require.Equal(t, result[0][0], "aaa" )
    require.Equal(t, result[0][1], "bbb" )

    msgs = []UniversalMsg{"aaa", "bbb", "ccc"}
    result = makeBatches(msgs, 3)
    require.Equal(t, len(result), 1 )
    require.Equal(t, len(result[0]), 3 )
    require.Equal(t, result[0][0], "aaa")
    require.Equal(t, result[0][1], "bbb")
    require.Equal(t, result[0][2], "ccc")

    msgs = []UniversalMsg{"aaa", "bbb", "ccc", "ddd"}
    result = makeBatches(msgs, 3)
    require.Equal(t, len(result), 2 )
    require.Equal(t, len(result[0]), 3 )
    require.Equal(t, len(result[1]), 1 )
    require.Equal(t, result[0][0], "aaa")
    require.Equal(t, result[0][1], "bbb")
    require.Equal(t, result[0][2], "ccc")
    require.Equal(t, result[1][0], "ddd")
}
