package schema

import (
	"encoding/json"

	"cuelang.org/go/cue"
)

// Decoder wraps a target and unmarshals JSON into it
type Decoder[T any] struct {
	Target *T
}

// UnmarshalJSON unmarhals the JSON data into the decoder target.
func (d *Decoder[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, d.Target)
}

// DecodeNew decodes the given [cue.Value] using the CUE/JSON unmarshaler into
// a new T. This is quite slow but currently the easiest way to get a CUE
// value into a Go struct. The type T can either use `json` tags on its
// fields, or implement [json.Unmarshaler].
func DecodeNew[T any](val cue.Value) (*T, error) {
	dec := Decoder[T]{new(T)}
	return dec.Target, val.Decode(&dec)
}

// Decode decodes the given [cue.Value] using the CUE/JSON unmarshaler into
// the given T pointer.
func Decode[T any](val cue.Value, iface *T) error {
	dec := Decoder[T]{iface}
	return val.Decode(&dec)
}
