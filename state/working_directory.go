package state

import "github.com/moby/buildkit/client/llb"

type WorkingDirectory struct {
	Directory string
}

func (wd *WorkingDirectory) StateOption() llb.StateOption {
	return func(state llb.State) llb.State {
		return state.Dir(wd.Directory)
	}
}
