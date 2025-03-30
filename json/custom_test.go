package json

import (
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

type customValue struct {
	x int64
}

var _ Appender = (*customValue)(nil)

func (v customValue) AppendJSON(b []byte, flags AppendFlags) ([]byte, error) {
	return Marshal(v.x)
}

var _ Parser = (*customValue)(nil)

func (v *customValue) ParseJSON(tok *Tokenizer, flags ParseFlags) error {
	if !tok.Next() {
		return io.ErrUnexpectedEOF
	}
	if kind := tok.Kind().Class(); kind != Num {
		return fmt.Errorf("unsupported token: %d", kind)
	}
	v.x = tok.Int()
	return nil
}

func TestCustomValue(t *testing.T) {
	src := (*customValue)(nil)
	buf, err := Marshal(src)
	require.NoError(t, err)
	require.Equal(t, []byte("null"), buf)

	values := []customValue{
		customValue{},
		customValue{42},
	}
	for _, src := range values {
		buf, err := Marshal(src)
		require.NoError(t, err)

		buf2, err := Marshal(&src)
		require.NoError(t, err)
		require.Equal(t, buf, buf2)

		var dest customValue
		err = Unmarshal(buf, &dest)
		require.NoError(t, err)
		require.Equal(t, src, dest)
	}
}
