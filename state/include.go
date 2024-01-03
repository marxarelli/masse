package state

import (
	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/masse/common"
)

type Include common.Include

func (inc *Include) SetLocalOption(info *llb.LocalInfo) {
	llb.IncludePatterns(globsToStrings(inc.Include)).SetLocalOption(info)
}

func (inc *Include) SetCopyOption(info *llb.CopyInfo) {
	info.IncludePatterns = globsToStrings(inc.Include)
}

func globsToStrings(globs []common.Glob) []string {
	strs := make([]string, len(globs))
	for i, glob := range globs {
		strs[i] = glob.String()
	}
	return strs
}
