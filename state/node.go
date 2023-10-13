package state

import (
	"fmt"
)

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
