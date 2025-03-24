package load

import (
	"cuelang.org/go/cue/load"
	"cuelang.org/go/mod/modconfig"
)

type Option func(dir string, cfg *load.Config, modcfg *modconfig.Config) error
