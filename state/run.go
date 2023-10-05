package state

type Run struct {
	Command   string `json:"run"`
	Arguments []string
	Options   []*RunOption
}

type RunOption struct {
	*Env
	*Host
	*CacheMount
	*SourceMount
	*TmpFSMount
	*ReadOnly
}

type SourceMount struct {
	Target string `json:"mount"`
	From   Chain
	Source string
}
