package state

import (
	"fmt"
	"strings"

	"github.com/moby/buildkit/client/llb"
)

type Diff struct {
	Upper Chain `json:"diff"`
}

func (diff *Diff) Description() string {
	descs := make([]string, len(diff.Upper))

	for i, state := range diff.Upper {
		descs[i] = state.Description()
	}

	return fmt.Sprintf(
		"Δ{ %s }",
		strings.Join(descs, ", "),
	)
}

func (diff *Diff) AnonymousChains() (chains []Chain, closed bool) {
	return []Chain{diff.Upper}, true
}

func (diff *Diff) Compile(lower llb.State, secondary ChainStates, constraints ...llb.ConstraintsOpt) (llb.State, error) {
	upper := lower
	var err error

	for _, state := range diff.Upper {
		cons := append(constraints, llb.WithCustomName(state.Description()))
		upper, err = state.Compile(upper, secondary, cons...)
		if err != nil {
			return lower, err
		}
	}

	return llb.Diff(lower, upper, constraints...), nil
}
