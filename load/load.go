package load

import (
	"bytes"
	"path/filepath"
	"text/template"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/mod/modconfig"
	"github.com/pkg/errors"
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

	modcfg := &modconfig.Config{}
	cfg := &load.Config{
		Dir:        dir,
		Package:    "main",
		ModuleRoot: ".",
	}

	rootCue, err := embedFS.ReadFile("root.cue")
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("module.cue").ParseFS(embedFS, "module.cue")
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, struct{ Version string }{schema.Version()})
	if err != nil {
		return nil, err
	}

	WithOverlayFiles(map[string][]byte{
		"root.cue":           rootCue,
		"cue.mod/module.cue": buf.Bytes(),
	})(dir, cfg, modcfg)

	for _, opt := range options {
		err = opt(dir, cfg, modcfg)
		if err != nil {
			return nil, err
		}
	}

	registry, err := NewRegistry(modcfg, cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "could not create new registry")
	}

	cfg.Registry = registry

	instances := load.Instances([]string{"."}, cfg)
	return instances[0], nil
}
