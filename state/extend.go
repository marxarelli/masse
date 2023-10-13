package state

import "github.com/moby/buildkit/client/llb"

type Extend struct {
	Extend ChainRef
}

func (ex *Extend) ChainRefs() []ChainRef {
	return []ChainRef{ex.Extend}
}

func (ex *Extend) Compile(_ llb.State, secondary ChainStates) (llb.State, error) {
	return secondary.Resolve(ex.Extend)
}
