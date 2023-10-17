package state

import (
	"fmt"
	"hash/fnv"
	"strconv"

	"github.com/dominikbraun/graph"
)

type NodeProperty = func(*graph.VertexProperties)

type Node struct {
	State     *State
	ChainRef  ChainRef
	Index     int
	Anonymous bool
}

func (n Node) Hash() string {
	if n.Index >= 0 {
		return fmt.Sprintf("%s[%d]", n.ChainRef, n.Index)
	}

	return string(n.ChainRef)
}

func (n Node) Properties() []NodeProperty {
	return []NodeProperty{
		n.Label(),
		n.Shape(),
		n.Style(),
		n.Color(),
	}
}

func (n Node) Description() string {
	return n.State.Description()
}

func (n Node) Label() NodeProperty {
	desc := strconv.Quote(n.Description())
	return graph.VertexAttribute(
		"label",
		desc[1:len(desc)-1],
	)
}

func (n Node) Color() NodeProperty {
	h := fnv.New32a()
	h.Write([]byte(n.ChainRef))

	return graph.VertexAttribute(
		"color",
		fmt.Sprintf("%d", (int(h.Sum32())%11)+1),
	)
}

func (n Node) Style() NodeProperty {
	style := "rounded"

	if n.Anonymous {
		style += ",dashed"
	} else {
		style += ",bold"
	}

	return graph.VertexAttribute("style", style)
}

func (n Node) Shape() NodeProperty {
	shape := "box"

	switch n.State.Kind() {
	case DiffKind:
		shape = "triangle"
	case MergeKind:
		shape = "invtriangle"
	}

	if n.ChainRef == ChainRef(".") {
		shape = "point"
	}

	return graph.VertexAttribute("shape", shape)
}
