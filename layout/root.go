package layout

import (
	"github.com/pkg/errors"
	"gitlab.wikimedia.org/dduvall/phyton/state"
)

type Root struct {
	Parameters Parameters   `json:"parameters"`
	Chains     state.Chains `json:"chains"`
	Layouts    Layouts      `json:"layouts"`
}

// LayoutGraph returns a new [state.Graph] for the given [Layouts] entry.
func (root *Root) LayoutGraph(name string) (*state.Graph, error) {
	layout, ok := root.Layouts[name]
	if !ok {
		return nil, errors.Errorf("unknown layout %q", name)
	}

	return layout.Graph(root.Chains)
}
