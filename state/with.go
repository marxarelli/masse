package state

import "github.com/moby/buildkit/client/llb"

type With struct {
	With Options `json:"withValue"`
}

func (with *With) Compile(state llb.State, _ ChainStates) (llb.State, error) {
	return with.With.StateOption()(state), nil
}
