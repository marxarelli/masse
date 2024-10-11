package v1

import (
	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/masse/common"
)

type Env struct {
	Env common.Env
}

func (env *Env) StateOption() llb.StateOption {
	return func(state llb.State) llb.State {
		for _, name := range env.Env.Sort() {
			state = state.AddEnv(name, env.Env[name])
		}

		return state
	}
}
