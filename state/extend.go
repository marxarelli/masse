package state

import "github.com/moby/buildkit/client/llb"

type Extend struct {
	Extend ChainRef
}

func (ex *Extend) ChainRefs() []ChainRef {
	return []ChainRef{ex.Extend}
}

func (ex *Extend) CompileSource(secondary ChainStates, _ ...llb.ConstraintsOpt) (llb.State, error) {
	return secondary.Resolve(ex.Extend)
}
