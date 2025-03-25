package v1

import (
	"cuelang.org/go/cue"
	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/masse/internal/lookup"
)

type FileActionKind string

const (
	CopyKind   FileActionKind = "copy"
	MkfileKind                = "mkfile"
	MkdirKind                 = "mkdir"
	RmKind                    = "rm"
)

func (c *compiler) compileFile(state llb.State, v cue.Value) (llb.State, error) {
	var fa *llb.FileAction

	err := lookup.EachOrValue(v, "file", func(filev cue.Value) error {
		var found bool
		var err error
		fa, found, err = lookup.WithDiscriminatorField(filev, func(kind FileActionKind) (*llb.FileAction, bool, error) {
			return c.compileDispatchFileKind(kind, state, fa, filev)
		})

		if err != nil {
			return err
		}

		if !found {
			return errorf(v, "unsupported file action")
		}

		return nil
	})

	if err != nil {
		return state, vError(v, err)
	}

	if fa == nil {
		return state, errorf(v, "file actions failed to compile")
	}

	options, err := lookup.DecodeOptions[Constraints](v)
	if err != nil {
		return state, vError(v, err)
	}

	return state.File(fa, options, c.opConstraints()), nil
}

func (c *compiler) compileDispatchFileKind(kind FileActionKind, state llb.State, fa *llb.FileAction, v cue.Value) (*llb.FileAction, bool, error) {
	var err error
	switch kind {
	case CopyKind:
		fa, err = c.compileCopy(state, fa, v)
	case MkfileKind, MkdirKind:
		fa, err = c.compileMkfile(kind, state, fa, v)
	case RmKind:
		fa, err = c.compileRm(state, fa, v)
	default:
		return fa, false, err
	}

	return fa, true, err
}
