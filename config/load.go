package config

import (
	"path/filepath"

	"cuelang.org/go/cue"
	"github.com/pkg/errors"
	"gitlab.wikimedia.org/dduvall/masse/load"
)

// Load loads the given configuration into a new *Root using a new CUE
// context. The path is used solely for location information in error
// reporting.
func Load(path string, data []byte, parameters map[string]string) (*Root, error) {
	val, err := LoadCUE(path, data, parameters)
	if err != nil {
		return nil, err
	}

	if val.Err() != nil {
		return nil, val.Err()
	}

	root := &Root{}
	return root, root.UnmarshalCUE(val)
}

// LoadCUE loads the given configuration and returns the root cue.Value.
func LoadCUE(path string, data []byte, parameters map[string]string) (cue.Value, error) {
	ctx := load.NewContext()

	basename := filepath.Base(path)
	dir := filepath.Dir(path)

	main, err := load.MainInstanceWith(
		dir,
		map[string][]byte{
			basename: data,
		},
	)

	if err != nil {
		return cue.Value{}, errors.Wrapf(err, "failed to parse configuration from %q", path)
	}

	root := ctx.BuildInstance(main)

	if parameters != nil {
		paramValues := make(map[string]cue.Value, len(parameters))
		for name, expr := range parameters {
			value := ctx.CompileString(expr)
			if err := value.Err(); err != nil {
				return root, errors.Wrapf(err, "failed to compile parameter %q", name)
			}

			paramValues[name] = value
		}

		root = root.FillPath(cue.ParsePath("parameters"), paramValues)
	}

	return root, nil
}
