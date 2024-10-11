package config

import (
	"path/filepath"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"github.com/pkg/errors"
	"gitlab.wikimedia.org/dduvall/masse/load"
)

// Load loads the given configuration into a new *Root using a new CUE
// context. The path is used solely for location information in error
// reporting.
func Load(path string, data []byte) (*Root, error) {
	val, err := LoadCUE(path, data)
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
func LoadCUE(path string, data []byte) (cue.Value, error) {
	ctx := cuecontext.New(cuecontext.EvaluatorVersion(cuecontext.EvalV2))

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

	return ctx.BuildInstance(main), nil
}
