package v1

import (
	"cuelang.org/go/cue"
	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/masse/internal/lookup"
)

func (c *compiler) compileDiff(lower llb.State, v cue.Value) (llb.State, error) {
	var err error
	upper := lower

	ref := lookup.Lookup(v, "from")
	if ref.Exists() && ref.IsConcrete() {
		lower, err = c.compileChainByRef(ref)
		if err != nil {
			return lower, err
		}
	}

	err = lookup.EachOrValue(v, "diff", func(opv cue.Value) error {
		if opv.IsNull() {
			return errorf(opv, "diff cannot have a null operation")
		}

		var err error
		upper, err = c.compileState(upper, opv)
		if err != nil {
			return vError(opv, err)
		}

		return nil
	})
	if err != nil {
		return lower, vError(v, err)
	}

	options, err := lookup.DecodeOptions[Constraints](v)
	if err != nil {
		return lower, vError(v, err)
	}

	return llb.Diff(lower, upper, options, c.opConstraints()), nil
}
