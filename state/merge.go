package state

import (
	"strings"

	"github.com/moby/buildkit/client/llb"
)

type Merge struct {
	Merge []ChainRef
}

func (mg *Merge) Description() string {
	refs := make([]string, len(mg.Merge))
	for i, ref := range mg.Merge {
		refs[i] = string(ref)
	}

	return strings.Join(refs, " ∪ ")
}

func (mg *Merge) ChainRefs() []ChainRef {
	return mg.Merge
}

func (mg *Merge) Compile(primary llb.State, secondary ChainStates, constraints ...llb.ConstraintsOpt) (llb.State, error) {
	return mg.compile(&primary, secondary, constraints...)
}

func (mg *Merge) CompileSource(secondary ChainStates, constraints ...llb.ConstraintsOpt) (llb.State, error) {
	return mg.compile(nil, secondary, constraints...)
}

func (mg *Merge) compile(primary *llb.State, secondary ChainStates, constraints ...llb.ConstraintsOpt) (llb.State, error) {
	states := make([]llb.State, len(mg.Merge))
	for i, ref := range mg.Merge {
		state, err := secondary.Resolve(ref)
		if err != nil {
			return llb.State{}, err
		}
		states[i] = state
	}

	if primary != nil {
		states = append([]llb.State{*primary}, states...)
	}

	return llb.Merge(states, constraints...), nil
}
