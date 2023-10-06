package state

type Chain []*State

func (chain Chain) Tail() (int, *State) {
	n := len(chain)
	if n > 0 {
		return n - 1, chain[n-1]
	}
	return 0, nil
}

type ChainRef string

type Chains map[ChainRef]Chain

type ChainReferencer interface {
	ChainRefs() []ChainRef
}
