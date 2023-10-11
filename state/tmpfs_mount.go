package state

import "github.com/moby/buildkit/client/llb"

type TmpFSMount struct {
	TmpFS string
	Size  uint64
}

func (tm *TmpFSMount) LLBRunOptions(_ ChainStates) ([]llb.RunOption, error) {
	return []llb.RunOption{
		llb.AddMount(
			tm.TmpFS,
			llb.Scratch(),
			llb.Tmpfs(llb.TmpfsSize(int64(tm.Size))),
		),
	}, nil
}
