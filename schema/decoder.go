package schema

import (
	"cuelang.org/go/cue"
)

// Unmarshaler defines an interface for types that can their own CUE value
// representations.
type Unmarshaler interface {
	UnmarshalCUE(cue.Value) error
}
