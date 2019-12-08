package golden

import (
	"errors"
	"testing"
)

func TestDefaultMarshallerErr(t *testing.T) {
	values := []interface{}{
		nil,
		struct{}{},
		errors.New(""),
	}
	for _, v := range values {
		_, err := DefaultMarshaller(v)
		if err != ErrUnsupportedType {
			t.Fatal("err != ErrUnsupportedType", v)
		}
	}
}

func TestDefaultMarshaller(t *testing.T) {
	values := []interface{}{
		[]byte{},
		[]byte("Hello, World!"),
		"Hello, World!",
		"",
		"\n\r\t",
	}
	for _, v := range values {
		_, err := DefaultMarshaller(v)
		if err != nil {
			t.Fatal("err != nil", v)
		}
	}
}
