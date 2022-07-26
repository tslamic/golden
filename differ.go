package golden

import "fmt"

// Differ returns a string showcasing the diff between u and v.
type Differ func(u, v []byte) string

func SimpleDiffer() Differ {
	return func(t, u []byte) string {
		return fmt.Sprintf("expected: \n\t%s\nactual:\n\t%s", t, u)
	}
}
