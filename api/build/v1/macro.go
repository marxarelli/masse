package v1

import "gopkg.in/yaml.v3"

var macroRegistry map[string]MacroResolver

func Register(name string, resolver MacroResolver) {
	// TODO mutex

	if macroRegistry == nil {
		macroRegistry = map[string]MacroResolver{}
	}

	macroRegistry[name] = resolver
}

// MacroResolver should return a struct value that implements [Macro].
type MacroResolver func() Macro

// Macro defines how macro types are expanded to [Op]
type Macro interface {
	yaml.Unmarshaler

	// Expand take the given macro arguments and returns an Op slice
	Expand(args ...any) ([]*Op, error)
}
