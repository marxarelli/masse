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

// Instances returns all CUE instances for our schema.
func Instances(ctx *cue.Context) ([]*cue.Instance, error) {
	cfg, err := LoaderConfig(ctx, "/")

	if err != nil {
		return nil, err
	}

	entries, err := fs.ReadDir(FS, ".")

	if err != nil {
		return nil, err
	}

	importPaths := []string{}
	for _, entry := range entries {
		if entry.IsDir() {
			importPaths = append(importPaths, filepath.Join(cuemod.Module(), "schema", entry.Name()))
		}
	}

	instances := cue.Build(load.Instances(importPaths, cfg))
	for _, instance := range instances {
		if instance.Err != nil {
			return nil, instance.Err
		}
	}

	return instances, nil
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