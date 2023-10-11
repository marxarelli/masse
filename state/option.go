package state

import "github.com/moby/buildkit/client/llb"

type LLBStateOption interface {
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

type Option struct {
	*Env
	*WorkingDirectory
}

func (opt *Option) StateOption() llb.StateOption {
	return func(state llb.State) llb.State {
		so, ok := oneof[LLBStateOption](opt)
		if ok {
			state = so.StateOption()(state)
		}

		return state
	}
}

func (opt *Option) LLBRunOptions(states ChainStates) ([]llb.RunOption, error) {
	return []llb.RunOption{opt.StateOption()}, nil
}
