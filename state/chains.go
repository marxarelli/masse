package state

import (
	"github.com/moby/buildkit/client/llb"
	"github.com/pkg/errors"
)

type Chain []*State

func (chain Chain) Tail() (int, *State) {
	n := len(chain)
	if n > 0 {
		return n - 1, chain[n-1]
	}
	return 0, nil
}

type ChainRef string

type Chains map[ChainRef]Chain

type ChainReferencer interface {
	ChainRefs() []ChainRef
}

type ChainDefiner interface {
	AnonymousChains() []Chain
}

type ChainStates map[ChainRef]llb.State

func (cs ChainStates) Resolve(ref ChainRef) (llb.State, error) {
	state, ok := cs[ref]
	if !ok {
		return state, errors.Errorf("state not found for chain %q", ref)
	}

	return state, nil
}
