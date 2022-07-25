package golden

import (
	"sync"
	"unicode"
)

// Transformer transforms and returns u and v, respectively.
type Transformer func(u, v []byte) ([]byte, []byte)

// StripWhitespace removes all whitespace characters.
var StripWhitespace Transformer = func(u, v []byte) ([]byte, []byte) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); u = strip(u) }()
	go func() { defer wg.Done(); v = strip(v) }()
	wg.Wait()
	return u, v
}

func strip(b []byte) []byte {
	buf := make([]byte, 0, len(b))
	for _, r := range b {
		if unicode.IsSpace(rune(r)) {
			continue
		}
		buf = append(buf, r)
	}
	return buf
}
