package v1

import (
	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/masse/common"
)

type Group common.Group

func (group *Group) SetCopyOption(info *llb.CopyInfo) {
	group.ChownOpt(info.ChownOpt).SetCopyOption(info)
}

func (group *Group) SetMkfileOption(info *llb.MkfileInfo) {
	group.ChownOpt(info.ChownOpt).SetMkfileOption(info)
}

func (group *Group) SetMkdirOption(info *llb.MkdirInfo) {
	group.ChownOpt(info.ChownOpt).SetMkdirOption(info)
}

func (group *Group) ChownOpt(other *llb.ChownOpt) *llb.ChownOpt {
	chown := &llb.ChownOpt{}

	if other != nil {
		chown.User = other.User
	}

	chown.Group = (&User{UID: group.GID, User: group.Group}).UserOpt()

	return chown
}
