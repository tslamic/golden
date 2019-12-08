package golden

import (
	"bytes"
	"errors"
	"flag"
	"io/ioutil"
	"os"
	"sync"
	"testing"
	"unicode"
)

var update = flag.Bool("update", false, "update golden files")

// Data represents the golden file attributes.
type Data struct {
	Path             string
	Perm             os.FileMode
	Marsh            Marshaller
	Diff             Differ
	Update           bool
	IgnoreWhitespace bool
}

// File creates a new golden file.
func File(path string, opts ...Option) *Data {
	d := &Data{
		Marsh:  DefaultMarshaller,
		Diff:   DefaultDiffer,
		Path:   path,
		Perm:   0644, // -rw-r--r--
		Update: *update,
	}
	for _, opt := range opts {
		opt(d)
	}
	return d
}

// Open opens a golden file. Closing of the file is the callers responsibility.
func (d *Data) Open() (*os.File, error) {
	return os.OpenFile(d.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, d.Perm)
}

// Possible errors when invoking the Eq func.
var (
	ErrNotEqual     = errors.New("not equal")
	ErrNoMarshaller = errors.New("no marshaller")
	ErrNoPath       = errors.New("no golden file path")
	ErrNoDiffer     = errors.New("no differ")
)

// Eq compares the value v to the contents of the golden file.
// If it's not equal, it returns the ErrNotEqual together with a diff string.
func (d *Data) Eq(v interface{}) (string, error) {
	if d.Path == "" {
		return "", ErrNoPath
	}
	if d.Marsh == nil {
		return "", ErrNoMarshaller
	}
	if d.Diff == nil {
		return "", ErrNoDiffer
	}
	m, err := d.Marsh(v)
	if err != nil {
		return "", err
	}
	if d.Update {
		if err = ioutil.WriteFile(d.Path, m, d.Perm); err != nil {
			return "", err
		}
	}
	f, err := ioutil.ReadFile(d.Path)
	if err != nil {
		return "", err
	}
	if d.IgnoreWhitespace {
		m, f = stripSpace(m, f)
	}
	if eq := bytes.Equal(m, f); eq {
		return "", nil
	}
	return d.Diff(m, f), ErrNotEqual
}

// Equals compares the value v to the contents of the golden file.
func (d *Data) Equals(t *testing.T, v interface{}) {
	diff, err := d.Eq(v)
	if err == nil {
		return
	}
	switch err {
	case ErrNotEqual:
		t.Fatalf("golden file does not match the param\n%s", diff)
	default:
		t.Fatal(err)
	}
}

func stripSpace(t, u []byte) ([]byte, []byte) {
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
