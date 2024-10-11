package v1

import (
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/client/llb/imagemetaresolver"
)

type CompilerOption interface {
	SetCompilerOption(*Config)
}

type compilerOption func(*Config)

func (f compilerOption) SetCompilerOption(cfg *Config) {
	f(cfg)
}

type Config struct {
	ImageMetaResolver llb.ImageMetaResolver
}

func newConfig(options []CompilerOption) *Config {
	cfg := &Config{
		ImageMetaResolver: imagemetaresolver.Default(),
	}

	for _, opt := range options {
		opt.SetCompilerOption(cfg)
	}

	return cfg
}

func WithImageMetaResolver(resolver llb.ImageMetaResolver) CompilerOption {
	return compilerOption(func(cfg *Config) {
		cfg.ImageMetaResolver = resolver
	})
}
