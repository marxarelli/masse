package v1

import (
	"context"
	"errors"

	"cuelang.org/go/cue"
	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/masse/internal/lookup"
)

func (c *compiler) compileContext(state llb.State, v cue.Value) (llb.State, error) {
	if c.config.NamedContextLoader == nil {
		return state, vError(v, errors.New("the compiler has not been configured to load named contexts"))
	}

	name, err := lookup.String(v, "context")
	if err != nil {
		return state, vError(v, err)
	}

	options, err := lookup.DecodeOptions[LocalOptions](v)
	if err != nil {
		return state, vError(v, err)
	}

	st, err := c.config.NamedContextLoader(c.ctx, name, options)
	if err != nil {
		return state, vError(v, err)
	}
	if st == nil {
		return state, errors.New("named context loader return nil state")
	}

	return *st, nil
}

func (c *compiler) compileMainContext(state llb.State, v cue.Value) (llb.State, error) {
	if c.config.MainContextLoader == nil {
		return state, vError(v, errors.New("the compiler has not been configured to load a main context"))
	}

	options, err := lookup.DecodeOptions[LocalOptions](v)
	if err != nil {
		return state, vError(v, err)
	}

	st, err := c.config.MainContextLoader(c.ctx, options)
	if err != nil {
		return state, vError(v, err)
	}
	if st == nil {
		return state, errors.New("main context loader return nil state")
	}

	return *st, nil
}

type NamedContextLoader func(ctx context.Context, name string, opts ...llb.LocalOption) (*llb.State, error)
type MainContextLoader func(ctx context.Context, opts ...llb.LocalOption) (*llb.State, error)

func WithNamedContextLoader(loader NamedContextLoader) CompilerOption {
	return compilerOption(func(cfg *Config) {
		cfg.NamedContextLoader = loader
	})
}

func WithMainContextLoader(loader MainContextLoader) CompilerOption {
	return compilerOption(func(cfg *Config) {
		cfg.MainContextLoader = loader
	})
}
