package state

type Merge struct {
	Merge []ChainRef
}

func (mg *Merge) ChainRefs() []ChainRef {
	return mg.Merge
}
