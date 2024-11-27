package v1

import (
	"cuelang.org/go/cue"
	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/masse/internal/lookup"
)

func (c *compiler) compileOptions(state llb.State, v cue.Value) (llb.State, error) {
	options, err := lookup.DecodeListOrSingle[Options](v, "options")
	if err != nil {
		return state, vError(v, err)
	}

	for _, option := range options {
		state = option.StateOption()(state)
	}

	return state, nil
}

type StateOption interface {
	StateOption() llb.StateOption
}

type Options []*Option

func (opts Options) StateOption() llb.StateOption {
	return func(state llb.State) llb.State {
		for _, opt := range opts {
			state = opt.StateOption()(state)
		}

		return state
	}
}

func (opts Options) SetRunOption(info *llb.ExecInfo) {
	llb.With(opts.StateOption()).SetRunOption(info)
}

type Option struct {
	*Env
	*WorkingDirectory
}

func (opt *Option) StateOption() llb.StateOption {
	return func(state llb.State) llb.State {
		so, ok := oneof[StateOption](opt)
		if ok {
			return so.StateOption()(state)
		}
		return state
	}
}

func (opt *Option) SetRunOption(info *llb.ExecInfo) {
	llb.With(opt.StateOption()).SetRunOption(info)
}
