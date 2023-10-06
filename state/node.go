package state

import (
	"fmt"

	"github.com/moby/buildkit/client/llb"
)

type Node struct {
	State *State
	Ref   ChainRef
	Index int
}

func (n Node) Location() string {
	if n.Index >= 0 {
		return fmt.Sprintf("%s[%d]", n.Ref, n.Index)
	}

	return string(n.Ref)
}

func (n Node) Compile(_inputs []llb.State) (llb.State, error) {
	// TODO implement
	return llb.Scratch(), nil
}
