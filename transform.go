package golden

import (
	"sync"
	"unicode"
)

// Transformer transforms and returns t and u, respectively.
type Transformer func(t, u []byte) ([]byte, []byte)

// StripWhitespace removes all whitespace characters.
var StripWhitespace Transformer = func(t, u []byte) ([]byte, []byte) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); t = strip(t) }()
	go func() { defer wg.Done(); u = strip(u) }()
	wg.Wait()
	return t, u
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
