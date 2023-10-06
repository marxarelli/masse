package state

import (
	"sort"

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
		inputs   map[string][]string
		compiled map[string]llb.State
	)

	// Note the PredecessorMap() is not deterministic as it uses nested maps.
	// Convert each second level map to a slice and sort it.
	pmap, err := g.PredecessorMap()
	if err != nil {
		return result, err
	}

	inputs = make(map[string][]string, len(pmap))

	for tail, headMap := range pmap {
		heads := make([]string, len(headMap))
		i := 0
		for head := range headMap {
			heads[i] = head
			i++
		}
		sort.Strings(heads)
		inputs[tail] = heads
	}

	size, _ := g.Order()
	compiled = make(map[string]llb.State, size)

	hashes, err := graph.StableTopologicalSort(g.Graph, func(x, y string) bool { return x < y })
	if err != nil {
		return result, err
	}

	for _, hash := range hashes {
		node, err := g.Vertex(hash)
		if err != nil {
			return result, errors.Wrap(err, "failed to solve graph")
		}

		// Get compiled inputs
		compiledInputs := make([]llb.State, len(inputs[hash]))
		for i, input := range inputs[hash] {
			compiledInputs[i] = compiled[input]
		}

		// Compile node
		state, err := node.Compile(compiledInputs)
		if err != nil {
			return result, err
		}

		compiled[hash] = state

		// Since the state graph has a universal (single) sink, we can simply
		// assume the last node to be compiled (in topo sort order) will be the
		// end result
		result = state
	}

	return result, nil
}
