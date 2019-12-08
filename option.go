package golden

import (
	"encoding/json"
	"encoding/xml"
)

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
		d.Add(StripWhitespace)
	}
)
