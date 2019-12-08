package golden

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"flag"
	"io/ioutil"
	"os"
	"sync"
	"testing"
	"unicode"

	"github.com/sergi/go-diff/diffmatchpatch"
)

var update = flag.Bool("update", false, "update golden files")

// Marshaller returns v encoded as []byte.
type Marshaller func(v interface{}) ([]byte, error)

// Differ returns a diff string between t and u.
type Differ func(t, u []byte) string

// Option can modify the golden file attributes.
type Option func(*Data)

// Convenience Options.
var (
	JSON Option = func(d *Data) {
		d.Marsh = json.Marshal
	}
	XML Option = func(d *Data) {
		d.Marsh = xml.Marshal
	}
	IgnoreWhitespace Option = func(d *Data) {
		d.IgnoreWhitespace = true
	}
)

// ErrUnsupportedType is returned when encoding values other than string or a []byte using the DefaultMarshaller.
var ErrUnsupportedType = errors.New("only []byte and string are supported by default, use a custom Marshaller, e.g. JSON")

// DefaultMarshaller can handle []byte or string type.
var DefaultMarshaller Marshaller = func(v interface{}) ([]byte, error) {
	switch v := v.(type) {
	case string:
		return []byte(v), nil
	case []byte:
		return v, nil
	default:
		return nil, ErrUnsupportedType
	}
}

var DefaultDiffer Differ = func(t, u []byte) string {
	p := diffmatchpatch.New()
	d := p.DiffMain(string(t), string(u), false)
	return p.DiffPrettyText(d)
}

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

func (d *Data) Write(b []byte) (int, error) {
	f, err := os.OpenFile(d.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, d.Perm)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.Write(b)
}

func (d *Data) Read(b []byte) (int, error) {
	f, err := os.Open(d.Path)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.Read(b)
}

// Possible errors when invoking the Eq func.
var (
	ErrNotEqual     = errors.New("not equal")
	ErrNoMarshaller = errors.New("no marshaller")
	ErrNoPath       = errors.New("no golden file path")
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
	cnt := 0
	for _, r := range b {
		if unicode.IsSpace(rune(r)) {
			continue
		}
		buf = append(buf, r)
		cnt++
	}
	return buf[:cnt]
}
