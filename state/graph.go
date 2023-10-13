package state

import (
	"fmt"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/dominikbraun/graph"
	"github.com/pkg/errors"
)

type Graph struct {
	graph.Graph[string, Node]
	chains      Chains
	addedChains sync.Map
	anonCounter atomic.Uint32
}

// NewGraph creates a new state DAG from the given [Chains] and terminal
// [*Merge].
func NewGraph(chains Chains, merge *Merge) (*Graph, error) {
	g := &Graph{
		Graph:       graph.New(nodeHash, graph.Directed(), graph.PreventCycles()),
		chains:      chains,
		addedChains: sync.Map{},
		anonCounter: atomic.Uint32{},
	}

	sink := Node{
		State:     &State{Merge: merge},
		ChainRef:  ChainRef("."),
		Index:     -1,
		Anonymous: false,
	}

	err := g.AddVertex(sink)
	if err != nil {
		return nil, err
	}

	for _, ref := range merge.ChainRefs() {
		err := g.AddChainEdgeByRef(ref, sink)
		if err != nil {
			return nil, err
		}
	}

	return g, nil
}

// AddChainEdgeByRef resolves a current chain by name only before calling
// [AddChainEdge].
func (g *Graph) AddChainEdgeByRef(ref ChainRef, node Node) error {
	chain, ok := g.chains[ref]
	if !ok {
		return errors.Errorf("unknown chain %q", ref)
	}

	return g.AddChainEdge(ref, chain, node, false)
}

// AddChainEdge adds a new edge (`Xn` -> `n`) where `Xn` is a node for the
// chain's sink and `n` is the given node. The entire chain is first added to
// the graph if it hasn't been already. If any state in the chain references
// other chains, [AddChainEdge] is called for each reference recursively. Note
// that this function assumes the given node has already been added to the
// graph.
func (g *Graph) AddChainEdge(ref ChainRef, chain Chain, node Node, anonymous bool) error {
	// Add entire chain first if it hasn't been added already
	_, exists := g.addedChains.LoadOrStore(ref, 1)
	if !exists {
		err := g.AddChain(ref, chain, anonymous)
		if err != nil {
			return errors.Wrapf(err, "error adding chain %q to graph", ref)
		}
	}

	// Then define an edge from the chain's sink to the given node
	i, sink := chain.Tail()
	if sink == nil {
		return errors.Errorf("%q chain sink is nil (chain is empty)", ref)
	}

	return g.AddEdge(Node{State: sink, ChainRef: ref, Index: i, Anonymous: anonymous}, node)
}

// AddAnonymousChainEdge resolves a new anonymous name for the given chain
// before calling [AddChainEdge] with it and the given node.
func (g *Graph) AddAnonymousChainEdge(chain Chain, node Node) error {
	return g.AddChainEdge(g.newAnonymous(node.Hash()), chain, node, true)
}

// AddChain adds a vertex for each state in the chain and an edge for each
// adjacent pair. It uses the given [ChainRef] and each index within the chain
// to uniquely identify each node. If any state in the chain references or
// defines anonymous chains, they are added via [AddChainEdgeByRef] or
// [AddAnonymousChainEdge] respectively.
func (g *Graph) AddChain(ref ChainRef, chain Chain, anonymous bool) error {
	var prev *Node
	for i, state := range chain {
		node := Node{State: state, ChainRef: ref, Index: i, Anonymous: anonymous}

		err := g.AddVertex(node)
		if err != nil {
			return err
		}

		// Does this state define anonymous chains? If so, add each other chain
		// first along with an edge from each chain sink to this state.
		for _, anonChain := range state.AnonymousChains() {
			err := g.AddAnonymousChainEdge(anonChain, node)
			if err != nil {
				return err
			}
		}

		// Does this state reference other chains? If so, add each other chain
		// first along with an edge from each chain sink to this state.
		for _, ref := range state.ChainRefs() {
			err := g.AddChainEdgeByRef(ref, node)
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

// InputMaps returns two maps derived from the [graph.Graph.PredecessorMap]
// that represents both the primary and secondary inputs of each graph node.
// The primary input of given node is the node that preceeds it along the same
// chain, while the secondary inputs are nodes that are referrenced explicitly
// by name.
func (g *Graph) InputMaps() (map[string]Node, map[string][]Node, error) {
	primary := map[string]Node{}
	secondary := map[string][]Node{}

	pmap, err := g.PredecessorMap()
	if err != nil {
		return primary, secondary, err
	}

	for tail, headMap := range pmap {
		tailNode, err := g.Vertex(tail)
		if err != nil {
			return primary, secondary, err
		}

		heads := make([]string, len(headMap))

		i := 0
		for head := range headMap {
			heads[i] = head
			i++
		}

		// Note the PredecessorMap() is not deterministic, so we need to sort
		sort.Strings(heads)

		for _, head := range heads {
			headNode, err := g.Vertex(head)
			if err != nil {
				return primary, secondary, err
			}

			// Separate primary from secondary inputs. The primary input is the
			// predecessor from the same chain
			if headNode.ChainRef == tailNode.ChainRef {
				primary[tail] = headNode
			} else {
				if _, ok := secondary[tail]; !ok {
					secondary[tail] = []Node{}
				}
				secondary[tail] = append(secondary[tail], headNode)
			}
		}
	}

	return primary, secondary, err
}

// newAnonymous increments the anonymous name counter and returns a name that
// can be used to identify an otherwise anonymous chain.
func (g *Graph) newAnonymous(scope string) ChainRef {
	return ChainRef(fmt.Sprintf("%s(anonymous%d)", scope, g.anonCounter.Add(uint32(1))))
}

func nodeHash(n Node) string {
	return n.Hash()
}
