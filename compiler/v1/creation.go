package v1

import (
	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/masse/common"
)

type Creation common.Creation

func (creation *Creation) SetCopyOption(info *llb.CopyInfo) {
	if creation.Ctime != nil {
		llb.WithCreatedTime(*creation.Ctime).SetCopyOption(info)
	}
}

func (creation *Creation) SetMkfileOption(info *llb.MkfileInfo) {
	if creation.Ctime != nil {
		llb.WithCreatedTime(*creation.Ctime).SetMkfileOption(info)
	}
}

func (creation *Creation) SetMkdirOption(info *llb.MkdirInfo) {
	if creation.Ctime != nil {
		llb.WithCreatedTime(*creation.Ctime).SetMkdirOption(info)
	}
}
