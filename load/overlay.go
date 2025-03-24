package load

import (
	"io/fs"
	"path"

	"cuelang.org/go/cue/load"
	"cuelang.org/go/mod/modconfig"
	"gitlab.wikimedia.org/dduvall/masse/util/readfs"
)

func WithFilesystem(subdir string, fsys fs.FS) Option {
	return func(dir string, cfg *load.Config, modcfg *modconfig.Config) error {
		files := map[string][]byte{}

		err := readfs.Read(fsys, ".", func(path string, data []byte) {
			files[path] = data
		})

		if err != nil {
			return err
		}

		return WithOverlayFiles(files)(path.Join(dir, subdir), cfg, modcfg)
	}
}

func WithOverlayFiles(files map[string][]byte) Option {
	return func(dir string, cfg *load.Config, _ *modconfig.Config) error {
		if cfg.Overlay == nil {
			cfg.Overlay = map[string]load.Source{}
		}

		for filename, data := range files {
			cfg.Overlay[path.Join(dir, filename)] = load.FromBytes(data)
		}

		return nil
	}
}
