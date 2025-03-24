package load

import (
	"errors"
	"os"
	"path/filepath"

	"cuelang.org/go/cue/load"
	"cuelang.org/go/mod/modconfig"
	"cuelang.org/go/mod/modfile"
	"gitlab.wikimedia.org/dduvall/masse/schema"
)

func WithModFileBytes(data []byte) Option {
	return func(dir string, cfg *load.Config, _ *modconfig.Config) error {
		cfg.Overlay[modFileIn(dir)] = load.FromBytes(data)
		return nil
	}
}

func WithModFile(modFile *modfile.File) Option {
	return func(dir string, cfg *load.Config, modcfg *modconfig.Config) error {
		data, err := modFile.Format()
		if err != nil {
			return err
		}

		return WithModFileBytes(data)(dir, cfg, modcfg)
	}
}

func WithDefaultEmbeddedModFile() Option {
	return WithModFile(schema.EmbeddedProjectModFile())
}

func WithNearestModFile() Option {
	return func(dir string, cfg *load.Config, modcfg *modconfig.Config) error {
		modFile, data, err := readNearestCueModFile(dir)
		if err != nil {
			return err
		}

		cfg.Overlay[modFile] = load.FromBytes(data)
		return nil
	}
}

func readNearestCueModFile(dir string) (string, []byte, error) {
	modFile := modFileIn(dir)
	data, err := os.ReadFile(modFile)

	if err == nil {
		return modFile, data, nil
	}

	if !errors.Is(err, os.ErrNotExist) {
		return "", data, err
	}

	subdir := filepath.Dir(dir)
	if subdir == dir {
		return "", data, nil
	}

	return readNearestCueModFile(subdir)
}

func modFileIn(dir string) string {
	return filepath.Join(dir, "cue.mod", "module.cue")
}
