package golden

import "github.com/sergi/go-diff/diffmatchpatch"

// Differ returns a diff string between t and u.
type Differ func(t, u []byte) string

// DefaultDiffer returns a colored diff string between t and u.
var DefaultDiffer Differ = func(t, u []byte) string {
	p := diffmatchpatch.New()
	d := p.DiffMain(string(t), string(u), false)
	return p.DiffPrettyText(d)
}
