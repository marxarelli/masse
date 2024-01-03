package schema

import (
	"io/fs"
	"path/filepath"

	"cuelang.org/go/cue/load"

	cuemod "gitlab.wikimedia.org/dduvall/masse/cue.mod"
)

// LoaderConfig returns a CUE [load.Config] that can load the embedded Masse
// schema definitions.
func LoaderConfig(root string) (*load.Config, error) {
	overlay := make(map[string]load.Source)

	err := fs.WalkDir(FS, ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !entry.Type().IsRegular() {
			return nil
		}

		cueData, err := FS.ReadFile(path)

		if err != nil {
			return err
		}

		overlay[pkgPath(root, path)] = load.FromBytes(cueData)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &load.Config{
		Dir:     root,
		Overlay: overlay,
	}, nil
}

func pkgPath(root, path string) string {
	return filepath.Join(
		root,
		"/cue.mod/pkg",
		cuemod.Module(),
		"schema",
		path,
	)
}
