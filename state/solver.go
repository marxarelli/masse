package state

import (
	"github.com/dominikbraun/graph"
	"github.com/moby/buildkit/client/llb"
	"github.com/pkg/errors"
	"gitlab.wikimedia.org/dduvall/phyton/common"
)

type Solver interface {
	Solve(g *Graph) (llb.State, error)
}

type platformSolver struct {
	platform    common.Platform
	constraints []llb.ConstraintsOpt
}

// NewPlatformSolver returns a new [Solver] for the given [common.Platform]
// and [llb.ConstraintsOpt].
func NewPlatformSolver(p common.Platform, constraints ...llb.ConstraintsOpt) Solver {
	sp := &Platform{p}
	return &platformSolver{
		platform:    p,
		constraints: append(constraints, sp.LLBConstraints()...),
	}
}

// Solve reduces a state graph to a singular [llb.State]. It walks the graph
// in topological sort order, compiling each node and passing the output to
// nodes on each outgoing edge.
func (s *platformSolver) Solve(g *Graph) (llb.State, error) {
	var (
		result   llb.State
		compiled map[string]llb.State
	)

	if g == nil {
		return result, errors.New("cannot solve nil graph")
	}

	primaryInputs, secondaryInputs, err := g.InputMaps()
	if err != nil {
		return result, err
	}

	size, err := g.Order()
	if err != nil {
		return result, err
	}

	compiled = make(map[string]llb.State, size)

	hashes, err := graph.TopologicalSort(g.Graph)
	if err != nil {
		return result, err
	}

	for _, hash := range hashes {
		node, err := g.Vertex(hash)
		if err != nil {
			return result, errors.Wrap(err, "failed to solve graph")
		}

		if node.Anonymous {
			continue
		}

		// Get compiled inputs
		var compiledPrimary *llb.State
		var compiledSecondary ChainStates

		if primary, ok := primaryInputs[hash]; ok {
			if cp, ok := compiled[primary.Hash()]; ok {
				compiledPrimary = &cp
			}
		}

		if secondary, ok := secondaryInputs[hash]; ok {
			compiledSecondary = make(ChainStates, len(secondary))
			for _, sec := range secondary {
				if compiledSec, ok := compiled[sec.Hash()]; ok {
					compiledSecondary[sec.ChainRef] = compiledSec
				}
			}
		}

		constraints := append(s.constraints, llb.WithCustomNamef(
			"[%s] %s", s.platform.ID(), node.Description(),
		))

		var state llb.State
		if compiledPrimary == nil {
			state, err = node.State.CompileSource(compiledSecondary, constraints...)
		} else {
			state, err = node.State.Compile(*compiledPrimary, compiledSecondary, constraints...)
		}

		if err != nil {
			return result, errors.Wrap(err, "failed to compile state")
		}

		compiled[hash] = state

		// Since the state graph has a universal sink, we can simply assume the
		// last node to be compiled (in topo sort order) will be the singular end
		// result
		result = state
	}

	return result, nil
}
