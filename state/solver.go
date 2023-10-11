package state

import (
	"github.com/dominikbraun/graph"
	"github.com/moby/buildkit/client/llb"
	"github.com/pkg/errors"
)

// Solver compiles all nodes in the graph to a single [llb.State].
type Solver func(*Graph) (llb.State, error)

// Solve reduces a state graph to a singular [llb.State]. It walks the graph
// in topological sort order, compiling each node and passing the output to
// nodes on each outgoing edge.
func Solve(g *Graph) (llb.State, error) {
	var (
		result   llb.State
		compiled map[string]llb.State
	)

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

		// Get compiled inputs
		var compiledPrimary llb.State
		var compiledSecondary ChainStates

		primary, ok := primaryInputs[hash]
		if ok {
			compiledPrimary, ok = compiled[primary.Hash()]
			if !ok {
				return result, errors.Errorf("primary input node %q not found in compile cache", primary.Hash())
			}
		} else {
			compiledPrimary = llb.Scratch()
		}

		secondary, ok := secondaryInputs[hash]
		if ok {
			compiledSecondary = make(ChainStates, len(secondary))
			for _, sec := range secondary {
				compiledSecondary[sec.ChainRef], ok = compiled[sec.Hash()]
				if !ok {
					return result, errors.Errorf("secondary input node %q not found in compile cache", sec.Hash())
				}
			}
		} else {
			compiledSecondary = ChainStates{}
		}

		// Compile node
		state, err := node.State.Compile(compiledPrimary, compiledSecondary)
		if err != nil {
			return result, err
		}

		compiled[hash] = state

		// Since the state graph has a universal sink, we can simply assume the
		// last node to be compiled (in topo sort order) will be the singular end
		// result
		result = state
	}

	return result, nil
}
