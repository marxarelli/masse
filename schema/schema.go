package schema

import (
	"embed"
	"io/fs"
	"path/filepath"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/load"

	cuemod "gitlab.wikimedia.org/dduvall/phyton/cue.mod"
)

//go:embed **/*.cue
var FS embed.FS

// LoaderConfig returns a CUE [load.Config] that can load the embedded Phyton
// schema definitions.
func LoaderConfig(ctx *cue.Context, root string) (*load.Config, error) {
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
		path,
	)
}
