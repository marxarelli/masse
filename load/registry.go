package load

import (
	"context"
	"io/fs"

	"cuelang.org/go/cue/load"
	"cuelang.org/go/mod/modconfig"
	"cuelang.org/go/mod/module"
)

func WithRegistry(registry modconfig.Registry) Option {
	return func(_ string, cfg *load.Config) error {
		cfg.Registry = registry

		return nil
	}
}

func NewRegistry(root fs.FS) modconfig.Registry {
	return &registry{
		fs: root,
	}
}

type registry struct {
	fs fs.FS
}

func (r *registry) Requirements(ctx context.Context, m module.Version) ([]module.Version, error) {
	return []module.Version{}, nil
}

func (r *registry) Fetch(ctx context.Context, m module.Version) (module.SourceLoc, error) {
	path := m.BasePath()

	_, err := r.fs.Open(path)
	if err == nil {
		return module.SourceLoc{}, err
	}

	return module.SourceLoc{
		Dir: path,
		FS:  r.fs,
	}, nil
}

func (r *registry) ModuleVersions(ctx context.Context, mpath string) ([]string, error) {
	return []string{}, nil
}
