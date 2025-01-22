package load

import (
	"path/filepath"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"

	"gitlab.wikimedia.org/dduvall/masse/schema"
)

// NewContext returns a new [*cue.Context] for evaluation of Masse CUE
// configuration.
func NewContext() *cue.Context {
	return cuecontext.New()
}

// MainInstance returns a CUE instance that unifies with a config.#Root
func MainInstance(dir string, options ...Option) (*build.Instance, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	cfg, err := schema.LoaderConfig(dir)
	if err != nil {
		return nil, err
	}

	cfg.Package = "main"

	rootCue, err := embedFS.ReadFile("root.cue")
	if err != nil {
		return nil, err
	}

	moduleCue, err := embedFS.ReadFile("module.cue")
	if err != nil {
		return nil, err
	}

	WithOverlayFiles(map[string][]byte{
		"root.cue":           rootCue,
		"cue.mod/module.cue": moduleCue,
	})(dir, cfg)

	for _, opt := range options {
		err = opt(dir, cfg)
		if err != nil {
			return nil, err
		}
	}

	instances := load.Instances([]string{"."}, cfg)
	return instances[0], nil
}
