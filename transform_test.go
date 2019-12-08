package golden

import (
	"bytes"
	"testing"
)

func TestStrip(t *testing.T) {
	values := map[string]string{
		" ":               "",
		"\n\t\v\r\f":      "",
		"abc\n\t123":      "abc123",
		"\vhello world\r": "helloworld",
	}
	for str, exp := range values {
		s := []byte(str)
		e := []byte(exp)
		if !bytes.Equal(strip(s), e) {
			t.Fatal("not equal", str, exp)
		}
	}
}
