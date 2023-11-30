package config

import (
	"github.com/pkg/errors"
	"gitlab.wikimedia.org/dduvall/phyton/state"
	"gitlab.wikimedia.org/dduvall/phyton/target"
)

type Root struct {
	Parameters Parameters     `json:"parameters"`
	Chains     state.Chains   `json:"chains"`
	Targets    target.Targets `json:"targets"`
}

// TargetGraph returns a new [state.Graph] for the given [Targets] entry.
func (root *Root) TargetGraph(name string) (*state.Graph, error) {
	target, ok := root.Targets[name]
	if !ok {
		return nil, errors.Errorf("unknown target %q", name)
	}

	return target.Graph(root.Chains)
}
