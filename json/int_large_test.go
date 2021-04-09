package json

import (
	"bytes"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringifyLargeInts(t *testing.T) {
	// From https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Number/MAX_SAFE_INTEGER
	require.Equal(t, 9007199254740991, jsMaxSafeInt)

	nn := []interface{}{
		int(jsMaxSafeInt + 2),
		uint(jsMaxSafeInt + 2),
		int64(jsMaxSafeInt + 2),
		uint64(jsMaxSafeInt + 2),
	}

	for _, n := range nn {
		var buf bytes.Buffer
		enc := NewEncoder(&buf)
		enc.SetStringifyLargeInts(true)

		{
			err := enc.Encode(n)
			require.Nil(t, err)

			got := buf.String()
			wanted := fmt.Sprintf(`"%d"`+"\n", n)
			require.Equal(t, wanted, got)
		}

		{
			buf.Reset()

			err := enc.Encode(Number(fmt.Sprint(n)))
			require.Nil(t, err)

			got := buf.String()
			wanted := fmt.Sprintf(`"%d"`+"\n", n)
			require.Equal(t, wanted, got)
		}
	}
}

func TestStringifyAndLargeInts(t *testing.T) {
	type Test struct {
		Foo int64 `json:",string"`
	}

	test := &Test{Foo: math.MaxInt64}

	b, err := Marshal(test)
	require.Nil(t, err)
	require.Equal(t, `{"Foo":"9223372036854775807"}`, string(b))

	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	enc.SetStringifyLargeInts(true)

	err = enc.Encode(test)
	require.Nil(t, err)
	require.Equal(t, `{"Foo":"9223372036854775807"}`+"\n", buf.String())
}
