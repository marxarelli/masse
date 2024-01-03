package state

import (
	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/masse/common"
)

type Group common.Group

func (group *Group) SetCopyOption(info *llb.CopyInfo) {
	chown := &llb.ChownOpt{}
	if info.ChownOpt == nil {
		chown = info.ChownOpt
	}

	chown.Group = (&User{UID: group.GID, User: group.Group}).UserOpt()
}
