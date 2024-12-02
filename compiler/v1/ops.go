package v1

import (
	"cuelang.org/go/cue"
	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/masse/internal/lookup"
)

func (c *compiler) compileOps(state llb.State, v cue.Value) (llb.State, error) {
	return state, lookup.EachOrValue(v, "ops", func(opv cue.Value) error {
		if opv.IsNull() {
			return errorf(opv, "ops cannot be null")
		}

		var err error
		state, err = c.compileState(state, opv)
		if err != nil {
			return vError(opv, err)
		}

		return nil
	})
}
