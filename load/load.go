package load

import (
	"path/filepath"

	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/load"

	"gitlab.wikimedia.org/dduvall/masse/schema"
)

// MainInstanceWith returns a CUE instance with no package that unifies with a
// config.#Root
func MainInstanceWith(dir string, files map[string][]byte) (*build.Instance, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	cfg, err := schema.LoaderConfig(dir)
	if err != nil {
		return nil, err
	}

	cfg.Package = "main"

	cueData, err := embedFS.ReadFile("root.cue")
	if err != nil {
		return nil, err
	}

	cfg.Overlay[filepath.Join(dir, "root.cue")] = load.FromBytes(cueData)

	cueData, err = embedFS.ReadFile("module.cue")
	if err != nil {
		return nil, err
	}

	cfg.Overlay[filepath.Join(dir, "cue.mod", "module.cue")] = load.FromBytes(cueData)

	for path, data := range files {
		cfg.Overlay[filepath.Join(dir, path)] = load.FromBytes(data)
	}

	instances := load.Instances([]string{"."}, cfg)
	return instances[0], nil
}
