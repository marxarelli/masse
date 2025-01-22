package load

import (
	"io/fs"
	"path"

	"cuelang.org/go/cue/load"
)

func WithFilesystem(fsys fs.FS) Option {
	return func(dir string, cfg *load.Config) error {
		files := map[string][]byte{}

		err := fs.WalkDir(fsys, ".", func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !entry.IsDir() {
				file, err := fsys.Open(path)
				defer file.Close()

				if err != nil {
					return err
				}

				var dt []byte
				_, err = file.Read(dt)
				if err != nil {
					return err
				}

				files[path] = dt
			}

			return nil
		})

		if err != nil {
			return err
		}

		return WithOverlayFiles(files)(dir, cfg)
	}
}

func WithOverlayFiles(files map[string][]byte) Option {
	return func(dir string, cfg *load.Config) error {
		for filename, data := range files {
			cfg.Overlay[path.Join(dir, filename)] = load.FromBytes(data)
		}

		return nil
	}
}
