package v1

import (
	"cuelang.org/go/cue"
	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/masse/internal/lookup"
)

func (c *compiler) compileExtend(state llb.State, v cue.Value) (llb.State, error) {
	ref, err := lookup.Existing(v, "extend")
	if err != nil {
		return state, vError(v, err)
	}

	return c.compileChainByRef(ref)
}
