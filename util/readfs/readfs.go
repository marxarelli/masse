package readfs

import "io/fs"

type Yield func(path string, data []byte)

func Read(fsys fs.FS, dir string, yield Yield) error {
	return fs.WalkDir(fsys, dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !entry.IsDir() {
			data, err := fs.ReadFile(fsys, path)
			if err != nil {
				return err
			}

			yield(path, data)
		}

		return nil
	})
}
