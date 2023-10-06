package state

import (
	"github.com/dominikbraun/graph"
	"github.com/pkg/errors"
)

type Graph struct {
	graph.Graph[string, Node]
	chains      Chains
	addedChains map[ChainRef]struct{}
}

// NewGraph creates a new state DAG from the given [Chains] and terminal
// [*Merge].
func NewGraph(chains Chains, merge *Merge) (*Graph, error) {
	g := &Graph{
		Graph:       graph.New(nodeHash, graph.Directed(), graph.PreventCycles()),
		chains:      chains,
		addedChains: map[ChainRef]struct{}{},
	}

	sink := Node{&State{Merge: merge}, ChainRef("."), -1}

	err := g.AddVertex(sink)
	if err != nil {
		return nil, err
	}

	for _, ref := range merge.ChainRefs() {
		err := g.AddChainEdge(ref, sink)
		if err != nil {
			return nil, err
		}
	}

	return g, nil
}

// AddChainEdge adds a new edge (`Xn` -> `n`) where `Xn` is a node for the
// chain's sink and `n` is the given node. The entire chain is first added
// to the graph if it hasn't been already. If any state in the chain
// references other chains, [AddChainEdge] is called for each reference
// recursively. Note that this function assumes the given node has already
// been added to the graph.
func (g *Graph) AddChainEdge(ref ChainRef, node Node) error {
	chain, ok := g.chains[ref]
	if !ok {
		return errors.Errorf("unknown chain %q", ref)
	}

	// Add entire chain first if it hasn't been added already
	_, exists := g.addedChains[ref]
	if !exists {
		err := g.AddChain(ref, chain)
		if err != nil {
			return errors.Wrapf(err, "error adding chain %q to graph", ref)
		}
	}
	g.addedChains[ref] = struct{}{}

	// Then define an edge from the chain's sink to the given node
	i, sink := chain.Tail()
	if sink == nil {
		return errors.Errorf("%q chain sink is nil (chain is empty)", ref)
	}

	return g.AddEdge(Node{sink, ref, i}, node)
}

// AddChain adds a vertex for each state in the chain and an edge for each
// adjacent pair. It uses the given [ChainRef] and each index within the chain
// to uniquely identify each node. If any state in the chain references other
// chains, [AddChainEdge] is called for each reference.
func (g *Graph) AddChain(ref ChainRef, chain Chain) error {
	var prev *Node
	for i, state := range chain {
		node := Node{State: state, Ref: ref, Index: i}

		err := g.AddVertex(node)
		if err != nil {
			return err
		}

		// Does this state reference other chains? If so, add each other chain
		// first along with an edge from each chain sink to this state.
		for _, ref := range state.ChainRefs() {
			err := g.AddChainEdge(ref, node)
			if err != nil {
				return err
			}
		}

		if prev != nil {
			err := g.AddEdge(*prev, node)
			if err != nil {
				return err
			}
		}

		prev = &node
	}

	return nil
}

// AddEdge defines an edge between vertices `x` and `y`.
func (g *Graph) AddEdge(x, y Node) error {
	return g.Graph.AddEdge(nodeHash(x), nodeHash(y))
}

func nodeHash(n Node) string {
	return n.Location()
}
