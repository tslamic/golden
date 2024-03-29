package golden

import (
	"bytes"
	"errors"
	"flag"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// update flag determines if a golden file should be updated.
var update = flag.Bool("update", false, "update golden files")

// Option allows golden file attributes changes.
// For example, to provide a custom Marshaller, you can do the following:
// 	gf := golden.File(name, func(attrs *Attrs) {
//    	attrs.Marshaller = func(v interface{}) ([]byte, error) {
//        	b := new(bytes.Buffer)
//        	err := gob.NewEncoder(b).Encode(v)
//        	return b.Bytes(), err
//    	}
// 	})
type Option func(*Attrs)

// Attrs represents the golden file attributes.
type Attrs struct {
	Path       string
	Flag       int         // Flag to use with os.OpenFile.
	Perm       os.FileMode // Perm to use with os.OpenFile and os.WriteFile.
	Update     bool
	Marshaller Marshaller
	Differ     Differ
	Transforms []Transformer
	ChunkSize  int64 // Byte size of the chunks used to read the golden file.
}

// Apply sets a new Transformer func, e.g.:
// 	gf := File(name).Apply(StripWhitespace)
func (d *Attrs) Apply(t ...Transformer) *Attrs {
	d.Transforms = append(d.Transforms, t...)
	return d
}

const (
	defaultFilePerm  = 0o644
	defaultChunkSize = 4096
)

// File creates a new golden file.
func File(path string, opts ...Option) *Attrs {
	d := &Attrs{
		Path:      path,
		Flag:      os.O_RDWR, //nolint:nosnakecase
		Perm:      defaultFilePerm,
		ChunkSize: defaultChunkSize,
		Update:    *update,
	}
	for _, opt := range opts {
		opt(d)
	}
	return d
}

var (
	ErrNoPath   = errors.New("missing golden file path")
	ErrNotEqual = errors.New("not equal")
)

// Eq compares the value v to the content of the golden file and returns
// ErrNotEqual together with a diff string generated by the Attr.Differ if they differ.
// It uses the Attr.Marshaller to encode v into a byte slice.
func (d *Attrs) Eq(v interface{}) (string, error) {
	if d.Path == "" {
		return "", ErrNoPath
	}
	if d.Marshaller == nil {
		// If no Marshaller is set, use the default one based on the file extension.
		ext := filepath.Ext(d.Path)
		m, err := marshallerFor(ext)
		if err != nil {
			return "", err
		}
		d.Marshaller = m
	}
	b, err := d.Marshaller(v)
	if err != nil {
		return "", err
	}
	if d.Update {
		err := os.WriteFile(d.Path, b, d.Perm)
		return "", err
	}
	f, err := os.OpenFile(d.Path, d.Flag, d.Perm)
	if err != nil {
		return "", err
	}
	defer f.Close()
	s, err := f.Stat()
	if err != nil {
		return "", err
	}
	fst, snd, err := diff(f, b, s.Size(), d.ChunkSize, d.Transforms)
	if err != nil {
		return "", err
	}
	if len(fst) == 0 && len(snd) == 0 {
		return "", nil
	}
	if d.Differ == nil {
		d.Differ = SimpleDiffer()
	}
	return d.Differ(fst, snd), ErrNotEqual
}

// Equals compares the value v to the contents of the golden file:
// 	expected := &greeter{Greeting: "Hello, World!"}
//	gf := golden.File("testdata/hello.json")
//	gf.Equals(t, expected)
func (d *Attrs) Equals(t *testing.T, v interface{}) {
	t.Helper()
	diff, err := d.Eq(v)
	if err != nil {
		if errors.Is(err, ErrNotEqual) {
			t.Fatalf("%s mismatch:\n%s", d.Path, diff)
		} else {
			t.Fatal(err)
		}
	}
}

func diff(r io.Reader, b []byte, size, chunk int64, t []Transformer) ([]byte, []byte, error) {
	buf := make([]byte, chunk)
	idx := int64(0)
	max := int64(len(b))
	for size > 0 {
		s := min(chunk, size)
		if _, err := r.Read(buf); err != nil {
			return nil, nil, err
		}
		u := buf[:s]
		v := b[idx:min(idx+s, max)]
		for _, tt := range t {
			u, v = tt(u, v)
		}
		if !bytes.Equal(u, v) {
			return u, v, nil
		}
		idx += s
		size -= s
	}
	return nil, nil, nil
}

func min(a, b int64) int64 {
	if a > b {
		return b
	}
	return a
}
