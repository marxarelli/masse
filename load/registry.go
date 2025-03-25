package load

import (
	"context"
	"io/fs"
	"net/http"

	"cuelang.org/go/cue/load"
	"cuelang.org/go/mod/modconfig"
	"cuelang.org/go/mod/module"
	"gitlab.wikimedia.org/dduvall/masse/schema"
)

func WithRegistryTransport(rt http.RoundTripper) Option {
	return func(_ string, _ *load.Config, modcfg *modconfig.Config) error {
		modcfg.Transport = rt
		return nil
	}
}

func NewRegistry(modcfg *modconfig.Config, cfg *load.Config) (modconfig.Registry, error) {
	reg, err := modconfig.NewRegistry(modcfg)
	if err != nil {
		return nil, err
	}

	mfs, err := schema.NewModuleFS(cfg)
	if err != nil {
		return nil, err
	}

	return &registry{
		Registry: reg,
		fs:       mfs,
	}, nil
}

type registry struct {
	modconfig.Registry
	fs fs.FS
}

func (reg *registry) Requirements(ctx context.Context, m module.Version) ([]module.Version, error) {
	if m.Equal(schema.ModuleVersion) {
		return schema.ModFile.DepVersions(), nil
	}

	return reg.Registry.Requirements(ctx, m)
}

func (reg *registry) Fetch(ctx context.Context, m module.Version) (module.SourceLoc, error) {
	if m.Equal(schema.ModuleVersion) {
		return module.SourceLoc{FS: reg.fs, Dir: "."}, nil
	}

	return reg.Registry.Fetch(ctx, m)
}
