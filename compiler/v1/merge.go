package v1

import (
	"cuelang.org/go/cue"
	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/masse/internal/lookup"
)

func (c *compiler) compileMerge(state llb.State, v cue.Value) (llb.State, error) {
	states := []llb.State{state}

	err := lookup.EachOrValue(v, "merge", func(refv cue.Value) error {
		chainState, err := c.compileChainByRef(refv)
		if err != nil {
			return vError(refv, err)
		}

		states = append(states, chainState)

		return nil
	})
	if err != nil {
		return state, vError(v, err)
	}

	return llb.Merge(states), nil
}
