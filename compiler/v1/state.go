package v1

import (
	"cuelang.org/go/cue"
	"github.com/moby/buildkit/client/llb"
)

type StateKind string

const (
	ScratchKind StateKind = "scratch"
	ExtendKind            = "extend"
	ImageKind             = "image"
	GitKind               = "git"
	LocalKind             = "local"
	RunKind               = "run"
)

func (c *compiler) compileState(state llb.State, v cue.Value) (llb.State, error) {
	iter, err := v.Fields(cue.Final(), cue.Concrete(true), cue.Optional(false))
	if err != nil {
		return state, err
	}

	for iter.Next() {
		sel := iter.Selector()

		if sel.LabelType() == cue.StringLabel {
			state, matched := c.compileDispatchKind(StateKind(sel.String()), state, v)

			if matched {
				return state, c.Error()
			}
		}
	}

	return state, c.addError(errorf(v, "unsupported operation"))
}

func (c *compiler) compileDispatchKind(kind StateKind, state llb.State, v cue.Value) (llb.State, bool) {
	var err error
	switch kind {
	case ScratchKind:
		state = llb.Scratch()
	case ExtendKind:
		state, err = c.compileExtend(state, v)
	case ImageKind:
		state, err = c.compileImage(state, v)
	case GitKind:
		state, err = c.compileGit(state, v)
	case LocalKind:
		state, err = c.compileLocal(state, v)
	case RunKind:
		state, err = c.compileRun(state, v)
	default:
		return state, false
	}

	c.addVError(v, err)

	return state, true
}
