package load

import (
	"os"
	"path/filepath"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"

	"gitlab.wikimedia.org/dduvall/phyton/schema"
)

// LoadPath loads a CUE file.
func LoadPath(path string) (cue.Value, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return cue.Value{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return cue.Value{}, err
	}

	return LoadBytes(data, path)
}

// LoadBytes loads CUE configuration. It requires a file path to initialize
// the [load.Config] overlay FS. However the file need not actually exist.
func LoadBytes(data []byte, path string) (cue.Value, error) {
	ctx := cuecontext.New()

	dir := filepath.Dir(path)
	cfg, err := schema.LoaderConfig(dir)
	if err != nil {
		return cue.Value{}, err
	}

	cfg.Overlay[path] = load.FromBytes(data)

	instances := load.Instances([]string{"."}, cfg)
	value := ctx.BuildInstance(instances[len(instances)-1])
	return value, value.Err()
}

// MainInstanceWith returns a CUE instance with no package that unifies with a
// layout.#Root
func MainInstanceWith(files map[string][]byte) (*build.Instance, error) {
	cfg, err := schema.LoaderConfig("/")
	if err != nil {
		return nil, err
	}

	cfg.Package = "_"

	cueData, err := embedFS.ReadFile("root.cue")
	if err != nil {
		return nil, err
	}

	cfg.Overlay["/root.cue"] = load.FromBytes(cueData)

	for path, data := range files {
		cfg.Overlay[path] = load.FromBytes(data)
	}

	instances := load.Instances([]string{"."}, cfg)
	return instances[0], nil
}
