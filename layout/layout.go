package layout

import (
	"github.com/moby/buildkit/client/llb"
	"github.com/pkg/errors"
	"gitlab.wikimedia.org/dduvall/phyton/common"
	"gitlab.wikimedia.org/dduvall/phyton/state"
)

type Layouts map[string]*Layout

type Layout struct {
	Comprises     []state.ChainRef  `json:"comprises"`
	Authors       []Author          `json:"authors"`
	Platforms     []common.Platform `json:"platformsValue"`
	Parameters    Parameters        `json:"parameters"`
	Configuration ImageConfig       `json:"configuration"`
}

// Graph returns a new [state.Graph] for the [Layout] and the given
// [state.Chains].
func (layout *Layout) Graph(chains state.Chains) (*state.Graph, error) {
	return state.NewGraph(chains, &state.Merge{Merge: layout.Comprises})
}

// Solvers returns a new [state.Solver] for each of the [Layout]'s platforms.
func (layout *Layout) Solvers(constraints ...llb.ConstraintsOpt) []state.Solver {
	solvers := make([]state.Solver, len(layout.Platforms))

	for i, platform := range layout.Platforms {
		solvers[i] = state.NewPlatformSolver(platform, constraints...)
	}

	return solvers
}

// ResolvePlatformSolver returns a single [state.Solver] for the given
// platform name. If the layout does not include the given platform, an error
// is returned.
func (layout *Layout) ResolvePlatformSolver(platformName string, constraints ...llb.ConstraintsOpt) (state.Solver, error) {
	platform, err := common.ParsePlatform(platformName)
	if err != nil {
		return nil, err
	}

	// Validate that the layout contains this platform
	found := false
	for _, p := range layout.Platforms {
		if p == platform {
			found = true
			break
		}
	}

	if !found {
		return nil, errors.Errorf("layout does not support platform %s", platformName)
	}

	return state.NewPlatformSolver(platform, constraints...), nil
}
