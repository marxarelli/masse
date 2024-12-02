package v1

import (
	"cuelang.org/go/cue"
	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/masse/internal/lookup"
)

type StateKind string

const (
	ScratchKind StateKind = "scratch"
	ExtendKind            = "extend"
	ImageKind             = "image"
	GitKind               = "git"
	LocalKind             = "local"
	RunKind               = "run"
	FileKind              = "file"
	MergeKind             = "merge"
	DiffKind              = "diff"
	OpsKind               = "ops"
	WithKind              = "with"
)

func (c *compiler) compileState(state llb.State, v cue.Value) (llb.State, error) {
	var found bool
	var err error

	state, found, err = lookup.WithDiscriminatorField(v, func(kind StateKind) (llb.State, bool, error) {
		return c.compileDispatchStateKind(kind, state, v)
	})

	if found {
		return state, err
	}

	return state, errorf(v, "unsupported operation: %s", v)
}

func (c *compiler) compileDispatchStateKind(kind StateKind, state llb.State, v cue.Value) (llb.State, bool, error) {
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
	case FileKind:
		state, err = c.compileFile(state, v)
	case MergeKind:
		state, err = c.compileMerge(state, v)
	case DiffKind:
		state, err = c.compileDiff(state, v)
	case OpsKind:
		state, err = c.compileOps(state, v)
	case WithKind:
		state, err = c.compileWith(state, v)
	default:
		return state, false, err
	}

	return state, true, err
}
