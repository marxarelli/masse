package load

import (
	"errors"
	"os"
	"path/filepath"

	"cuelang.org/go/cue/load"
	"cuelang.org/go/mod/modconfig"
)

func WithNearestModFile() Option {
	return func(dir string, cfg *load.Config, modcfg *modconfig.Config) error {
		modFile, data, err := readCueModFile(dir)
		if err != nil {
			return err
		}

		cfg.Overlay[modFile] = load.FromBytes(data)
		return nil
	}
}

func readCueModFile(dir string) (string, []byte, error) {
	modFile := filepath.Join(dir, "cue.mod", "module.cue")
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

	return readCueModFile(subdir)
}
