package golden

import (
	"errors"
)

// Marshaller returns v encoded as []byte.
type Marshaller func(v interface{}) ([]byte, error)

// ErrUnsupportedType is returned when encoding values other than string or a []byte using the DefaultMarshaller.
var ErrUnsupportedType = errors.New("only []byte and string are supported by default, use a custom Marshaller, e.g. JSON")

// DefaultMarshaller can handle []byte or string type.
var DefaultMarshaller Marshaller = func(v interface{}) ([]byte, error) {
	if v == nil {
		return nil, ErrUnsupportedType
	}
	switch v := v.(type) {
	case string:
		return []byte(v), nil
	case []byte:
		return v, nil
	default:
		return nil, ErrUnsupportedType
	}
}
