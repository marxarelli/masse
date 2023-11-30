package config

import (
	"path/filepath"

	"cuelang.org/go/cue/cuecontext"
	"github.com/pkg/errors"
	"gitlab.wikimedia.org/dduvall/phyton/load"
	"gitlab.wikimedia.org/dduvall/phyton/schema"
)

// Load loads the given target configuration into a new *Root using a new CUE
// context. The path is used solely for location information in error
// reporting.
func Load(path string, data []byte) (*Root, error) {
	basename := filepath.Join("/", filepath.Base(path))

	main, err := load.MainInstanceWith(map[string][]byte{
		basename: data,
	})

	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse layout from %q", path)
	}

	ctx := cuecontext.New()
	val := ctx.BuildInstance(main)
	if val.Err() != nil {
		return nil, val.Err()
	}

	return schema.DecodeNew[Root](val)
}
