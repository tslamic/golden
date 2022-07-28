package golden_test

import (
	"bytes"
	"testing"

	"github.com/tslamic/golden/v2"
)

func TestStripWhitespace(t *testing.T) {
	values := [][2]string{
		{" ", "\n\t\v\r\f"},
		{"abc\n\t123", "abc123"},
		{"\vhello world\r", "\n\thello\v\r\fworld"},
	}
	for _, str := range values {
		x := []byte(str[0])
		y := []byte(str[1])
		u, v := golden.StripWhitespace(x, y)
		if !bytes.Equal(u, v) {
			t.Fatalf("not equal: %s, %s", x, y)
		}
	}
}
