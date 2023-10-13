package state

import (
	"github.com/moby/buildkit/client/llb"
)

type Diff struct {
	Upper Chain `json:"diff"`
}

func (diff *Diff) AnonymousChains() []Chain {
	return []Chain{diff.Upper}
}

func (diff *Diff) Compile(lower llb.State, secondary ChainStates) (llb.State, error) {
	upper := lower
	var err error

	for _, state := range diff.Upper {
		upper, err = state.Compile(upper, secondary)
		if err != nil {
			return lower, err
		}
	}

	return llb.Diff(lower, upper), nil
}
