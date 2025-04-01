package v1

import (
	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/masse/common"
)

type Env struct {
	Env             common.Env `json:"env"`
	ExpandVariables bool       `json:"expandVariables"`
	c               *compiler
}

func (env *Env) SetCompiler(c *compiler) {
	env.c = c
}

func (env *Env) StateOption() llb.StateOption {
	return func(state llb.State) llb.State {
		for _, name := range env.Env.Sort() {
			value := env.Env[name]

			if env.ExpandVariables {
				value = env.c.Expand(state, value)
			}

			state = state.AddEnv(name, value)
		}

		return state
	}
}
