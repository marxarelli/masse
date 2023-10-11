package state

import (
	"github.com/moby/buildkit/client/llb"
	"github.com/pkg/errors"
)

type Diff struct {
	Upper Chain `json:"diff"`
}

func (diff *Diff) AnonymousChains() []Chain {
	return []Chain{diff.Upper}
}

func (diff *Diff) Compile(lower llb.State, secondary ChainStates) (llb.State, error) {
	if len(secondary) != 1 {
		return lower, errors.Errorf("diff should have exact one secondary input")
	}

	var upper llb.State
	for _, state := range secondary {
		upper = state
	}

	return llb.Diff(lower, upper), nil
}
