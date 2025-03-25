package schema

import (
	"embed"
	"path"

	"cuelang.org/go/cue/load"
	"cuelang.org/go/mod/module"
	"gitlab.wikimedia.org/dduvall/masse/util/readfs"
)

const (
	masseEmbeddedModDir = "cue.mod/embed/masse"
)

func NewModuleFS(cfg *load.Config) (module.OSRootFS, error) {
	root := path.Join(cfg.Dir, masseEmbeddedModDir)

	// the cue loader will not actually load the file from our embedded FS, so
	// we must add the files to the load.Config.Overlay. We'll not put them
	// under `cue.mod/{usr,pkg,gen}` however, as this would cause ambiguous
	// import errors.
	if cfg.Overlay == nil {
		cfg.Overlay = map[string]load.Source{}
	}

	err := readfs.Read(FS, ".", func(file string, data []byte) {
		cfg.Overlay[path.Join(root, file)] = load.FromBytes(data)
	})

	return &moduleFS{FS: FS, root: root}, err
}

type moduleFS struct {
	embed.FS
	root string
}

func (mfs *moduleFS) OSRoot() string {
	return mfs.root
}
