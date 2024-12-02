package v1

import (
	"strings"

	"cuelang.org/go/cue"
	"github.com/moby/buildkit/client/llb"
	"github.com/pkg/errors"
)

func (c *compiler) compileChain(v cue.Value) (llb.State, error) {
	var err error
	state := llb.Scratch()

	err = v.Null()
	if err == nil {
		return state, nil
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

	return state, err
}

func (c *compiler) compileChainByRef(refv cue.Value) (llb.State, error) {
	state := llb.NewState(nil)

	ref, err := refv.String()
	if err != nil {
		return state, c.addVError(refv, err)
	}

	for _, r := range c.refStack {
		if ref == r {
			return state, c.addVError(
				refv,
				errors.Errorf("chain ref cycle detected: %s -> %s", strings.Join(c.refStack, " -> "), ref),
			)
		}
	}

	cc, ok := c.chainCompilers[ref]
	if !ok {
		return state, c.addVError(refv, errors.Errorf("unknown chain %q", ref))
	}

	result := cc(c.withRefOnStack(ref))
	return result.state, c.addVError(refv, result.err)
}
