package v1

import (
	"strings"

	"cuelang.org/go/cue"
	"github.com/moby/buildkit/client/llb"
	"github.com/pkg/errors"
	"gitlab.wikimedia.org/dduvall/masse/internal/lookup"
)

func (c *compiler) compileChain(v cue.Value) (llb.State, error) {
	v = lookup.Dereference(v)
	ref := lookup.NormalizeReference(v)

	res, err, _ := c.group.Do(ref, func() (any, error) {
		var err error
		state := llb.Scratch()

		err = v.Null()
		if err == nil {
			return state, nil
		}

		if v.Kind() != cue.ListKind {
			return c.compileState(state, v)
		}

		list, err := v.List()
		if err != nil {
			return state, vError(v, err)
		}

		for list.Next() {
			state, err = c.compileState(state, list.Value())

			if err != nil {
				return state, err
			}
		}

		c.states.Store(ref, state)

		return state, err
	})

	return res.(llb.State), err
}

func (c *compiler) compileChainByRef(chainOrRefv cue.Value) (llb.State, error) {
	state := llb.NewState(nil)

	ref, err := chainOrRefv.String()
	if err != nil {
		// if value is not a string, assuming it's an actual #Chain
		return c.compileChain(chainOrRefv)
	}

	for _, r := range c.refStack {
		if ref == r {
			return state, c.addVError(
				chainOrRefv,
				errors.Errorf("chain ref cycle detected: %s -> %s", strings.Join(c.refStack, " -> "), ref),
			)
		}
	}

	cc, ok := c.chainCompilers[ref]
	if !ok {
		return state, c.addVError(chainOrRefv, errors.Errorf("unknown chain %q", ref))
	}

	result := cc(c.withRefOnStack(ref))
	return result.state, c.addVError(chainOrRefv, result.err)
}
