package golden

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"strings"
)

var (
	ErrUnknownExtension = errors.New("unknown extension, please provide a custom Marshaller")
	ErrInvalidType      = errors.New("arg has an invalid type")
)

// Marshaller encodes v into a byte slice.
type Marshaller func(v interface{}) ([]byte, error)

// ByteMarshal assumes v will either be a string or a byte slice
// and marshals it accordingly.
func ByteMarshal(v interface{}) ([]byte, error) {
	switch v := v.(type) {
	case string:
		return []byte(v), nil
	case []byte:
		return v, nil
	default:
		return nil, ErrInvalidType
	}
}

func marshallerFor(extension string) (Marshaller, error) {
	ext := strings.ToLower(extension)
	switch ext {
	case ".json":
		return json.Marshal, nil
	case ".xml":
		return xml.Marshal, nil
	case ".txt":
		return ByteMarshal, nil
	default:
		return nil, ErrUnknownExtension
	}
}
