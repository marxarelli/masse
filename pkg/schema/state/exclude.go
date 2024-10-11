package state

import (
	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/masse/common"
)

type Exclude common.Exclude

func (ex *Exclude) SetLocalOption(info *llb.LocalInfo) {
	llb.ExcludePatterns(globsToStrings(ex.Exclude)).SetLocalOption(info)
}

func (ex *Exclude) SetCopyOption(info *llb.CopyInfo) {
	info.ExcludePatterns = globsToStrings(ex.Exclude)
}
