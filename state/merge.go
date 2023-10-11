package state

import "github.com/moby/buildkit/client/llb"

type Merge struct {
	Merge []ChainRef
}

func (mg *Merge) ChainRefs() []ChainRef {
	return mg.Merge
}

func (mg *Merge) Compile(primary llb.State, secondary ChainStates) (llb.State, error) {
	states := make([]llb.State, len(mg.Merge))
	for i, ref := range mg.Merge {
		state, err := secondary.Resolve(ref)
		if err != nil {
			return primary, err
		}
		states[i] = state
	}

	return llb.Merge(append([]llb.State{primary}, states...)), nil
}
