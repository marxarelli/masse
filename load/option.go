package load

import (
	"cuelang.org/go/cue/load"
)

type Option func(dir string, cfg *load.Config) error
