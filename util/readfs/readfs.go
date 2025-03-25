package readfs

import "io/fs"

type YieldRead func(path string, data []byte)

// Read walks the entire given fs.FS starting at dir, reads each file and
// passes the path and content to the YieldRead function.
func Read(fsys fs.FS, dir string, yield YieldRead) error {
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
