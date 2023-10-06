package state

type Run struct {
	Command   string `json:"run"`
	Arguments []string
	Options   []*RunOption
}

type RunOption struct {
	*Host
	*CacheMount
	*SourceMount
	*TmpFSMount
	*ReadOnly
	*Option
}

type SourceMount struct {
	Target string `json:"mount"`
	From   ChainRef
	Source string
}

func (run *Run) ChainRefs() []ChainRef {
	refs := []ChainRef{}

	for _, op := range run.Options {
		if op.SourceMount != nil {
			refs = append(refs, op.SourceMount.From)
		}
	}

	return refs
}
