package state

import (
	"fmt"

	"github.com/moby/buildkit/client/llb"
)

type With struct {
	With Options `json:"withValue"`
}

func (with *With) Description() string {
	return fmt.Sprintf("(%d options)", len(with.With))
}

func (with *With) Compile(state llb.State, _ ChainStates, _ ...llb.ConstraintsOpt) (llb.State, error) {
	return with.With.StateOption()(state), nil
}
