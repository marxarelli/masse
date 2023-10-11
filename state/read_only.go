package state

import "github.com/moby/buildkit/client/llb"

type ReadOnly struct {
	ReadOnly bool
}

func (ro *ReadOnly) LLBRunOptions(states ChainStates) ([]llb.RunOption, error) {
	ropts := []llb.RunOption{}

	if ro.ReadOnly {
		ropts = append(ropts, llb.ReadonlyRootFS())
	}

	return ropts, nil
}
