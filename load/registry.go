package load

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"

	"cuelang.org/go/cue/load"
	"cuelang.org/go/mod/modconfig"
	"cuelang.org/go/mod/modfile"
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
		reg: reg,
		fs:  mfs,
	}, nil
}

type registry struct {
	reg modconfig.Registry
	fs  fs.FS
}

var (
	_ oldRegistry = (*registry)(nil)
	_ newRegistry = (*registry)(nil)
)

// TODO remove this when we can assume we don't need to deal with
// the old Registry interface any more.
type oldRegistry interface {
	Fetch(ctx context.Context, m module.Version) (module.SourceLoc, error)
	Requirements(ctx context.Context, m module.Version) ([]module.Version, error)
}

type newRegistry interface {
	Fetch(ctx context.Context, m module.Version) (module.SourceLoc, error)
	ModFile(ctx context.Context, m module.Version) (*modfile.File, error)
}

func (reg *registry) Requirements(ctx context.Context, m module.Version) ([]module.Version, error) {
	if m.Equal(schema.ModuleVersion) {
		return schema.ModFile.DepVersions(), nil
	}
	switch r := reg.reg.(type) {
	case oldRegistry:
		return r.Requirements(ctx, m)
	case newRegistry:
		mf, err := r.ModFile(ctx, m)
		if err != nil {
			return nil, err
		}
		return mf.DepVersions(), nil
	default:
		panic("unreachable")
	}
}

func (reg *registry) Fetch(ctx context.Context, m module.Version) (module.SourceLoc, error) {
	if m.Equal(schema.ModuleVersion) {
		return module.SourceLoc{FS: reg.fs, Dir: "."}, nil
	}

	return reg.reg.Fetch(ctx, m)
}

func (reg *registry) ModFile(ctx context.Context, m module.Version) (*modfile.File, error) {
	if m.Equal(schema.ModuleVersion) {
		return schema.ModFile, nil
	}
	if r, ok := reg.reg.(newRegistry); ok {
		return r.ModFile(ctx, m)
	}
	return nil, fmt.Errorf("ModFile method not supported on old modconfig.Registry interface")
}

func (reg *registry) ModuleVersions(ctx context.Context, mpath string) ([]string, error) {
	return reg.reg.ModuleVersions(ctx, mpath)
}
