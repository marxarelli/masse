package v1

import (
	"cuelang.org/go/cue"
	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/masse/internal/lookup"
)

func (c *compiler) compileWith(state llb.State, v cue.Value) (llb.State, error) {
	options, err := lookup.DecodeListOrSingle[Options](v, "with")
	if err != nil {
		return state, vError(v, err)
	}

	for _, option := range options {
		state = option.StateOption()(state)
	}

	return state, nil
}
