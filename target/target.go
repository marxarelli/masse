package target

import (
	"github.com/moby/buildkit/client/llb"
	"github.com/pkg/errors"
	"gitlab.wikimedia.org/dduvall/phyton/common"
	"gitlab.wikimedia.org/dduvall/phyton/state"
)

type Targets map[string]*Target

type Target struct {
	Build     state.ChainRef    `json:"build"`
	Platforms []common.Platform `json:"platformsValue"`
	Runtime   Runtime           `json:"runtime"`
	Labels    map[string]string `json:"labels"`
}

// Graph returns a new [state.Graph] for the [Target] and the given
// [state.Chains].
func (target *Target) Graph(chains state.Chains) (*state.Graph, error) {
	return state.NewGraph(chains, target.Build)
}

// Solvers returns a new [state.Solver] for each of the [Target]'s platforms.
func (target *Target) Solvers(constraints ...llb.ConstraintsOpt) []state.Solver {
	solvers := make([]state.Solver, len(target.Platforms))

	for i, platform := range target.Platforms {
		solvers[i] = state.NewPlatformSolver(platform, constraints...)
	}

	return solvers
}

// ResolvePlatformSolver returns a single [state.Solver] for the given
// platform name. If the target does not include the given platform, an error
// is returned.
func (target *Target) ResolvePlatformSolver(platformName string, constraints ...llb.ConstraintsOpt) (state.Solver, error) {
	platform, err := common.ParsePlatform(platformName)
	if err != nil {
		return nil, err
	}

	// Validate that the target contains this platform
	found := false
	for _, p := range target.Platforms {
		if p == platform {
			found = true
			break
		}
	}

	if !found {
		return nil, errors.Errorf("target does not support platform %s", platformName)
	}

	return state.NewPlatformSolver(platform, constraints...), nil
}
