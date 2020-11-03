package json

import (
	"strconv"
	"unsafe"
)

// jsMaxSafeInt is the maximum safe integer in JavaScript.
//
// See https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Number/MAX_SAFE_INTEGER
const jsMaxSafeInt = 1<<53 - 1

func (e encoder) appendInt(b []byte, n int64) []byte {
	if e.flags&stringifyInts != 0 || (e.flags&StringifyLargeInts != 0 && n > jsMaxSafeInt) {
		b = append(b, '"')
		b = strconv.AppendInt(b, n, 10)
		b = append(b, '"')
		return b
	}
	return strconv.AppendInt(b, n, 10)
}

func (e encoder) appendUint(b []byte, n uint64) []byte {
	if e.flags&stringifyInts != 0 || (e.flags&StringifyLargeInts != 0 && n > jsMaxSafeInt) {
		b = append(b, '"')
		b = strconv.AppendUint(b, n, 10)
		b = append(b, '"')
		return b
	}
	return strconv.AppendUint(b, n, 10)
}

func constructIntToStringEncodeFunc(encode encodeFunc) encodeFunc {
	return func(e encoder, b []byte, p unsafe.Pointer) ([]byte, error) {
		e.flags |= stringifyInts

		b, err := encode(e, b, p)
		if err != nil {
			return b, err
		}
		return b, nil
	}
}
