package base62

import (
	"fmt"
	"strings"

	"github.com/douglasmakey/ursho/encoding"
)

// All characters
const (
	alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	length   = int64(len(alphabet))
)

func New() encoding.Coder {
	return &coder{}
}

type coder struct{}

// Encode number to base62.
func (*coder) Encode(n int64) string {
	if n == 0 {
		return string(alphabet[0])
	}

	s := ""
	for ; n > 0; n = n / length {
		s = string(alphabet[n%length]) + s
	}
	return s
}

// Decode converts a base62 token to int.
func (*coder) Decode(key string) (int64, error) {
	var n int64
	for _, c := range []byte(key) {
		i := strings.IndexByte(alphabet, c)
		if i < 0 {
			return 0, fmt.Errorf("unexpected character %c in base62 literal", c)
		}
		n = length*n + int64(i)
	}
	return n, nil
}
