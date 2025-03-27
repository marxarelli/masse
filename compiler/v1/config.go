package v1

import (
	"context"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/client/llb/imagemetaresolver"
	"gitlab.wikimedia.org/dduvall/masse/common"
)

type CompilerOption interface {
	SetCompilerOption(*Config)
}

type compilerOption func(*Config)

func (f compilerOption) SetCompilerOption(cfg *Config) {
	f(cfg)
}

type Config struct {
	ImageMetaResolver  llb.ImageMetaResolver
	InitialContext     context.Context
	Platform           *Platform
	IgnoreCache        bool
	NamedContextLoader NamedContextLoader
	MainContextLoader  MainContextLoader
}

func (cfg *Config) OpConstraints() Constraints {
	return Constraints{
		&Constraint{
			IgnoreCache: &IgnoreCache{
				IgnoreCache: cfg.IgnoreCache,
			},
		},
	}
}

func (cfg *Config) SourceConstraints() Constraints {
	return Constraints{
		&Constraint{
			Platform: cfg.Platform,
		},
	}
}

func newConfig(options []CompilerOption) *Config {
	cfg := &Config{
		InitialContext:    context.Background(),
		ImageMetaResolver: imagemetaresolver.Default(),
		Platform:          defaultPlatform,
	}

	for _, opt := range options {
		opt.SetCompilerOption(cfg)
	}

	return cfg
}

func WithContext(ctx context.Context) CompilerOption {
	return compilerOption(func(cfg *Config) {
		cfg.InitialContext = ctx
	})
}

func WithIgnoreCache(ignore bool) CompilerOption {
	return compilerOption(func(cfg *Config) {
		cfg.IgnoreCache = ignore
	})
}

func WithImageMetaResolver(resolver llb.ImageMetaResolver) CompilerOption {
	return compilerOption(func(cfg *Config) {
		cfg.ImageMetaResolver = resolver
	})
}

func WithPlatform(platform common.Platform) CompilerOption {
	return compilerOption(func(cfg *Config) {
		cfg.Platform = &Platform{Platform: platform}
	})
}
