package state

import (
	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/phyton/common"
)

type Creation common.Creation

func (creation *Creation) SetCopyOption(info *llb.CopyInfo) {
	info.CreatedTime = creation.Ctime
}