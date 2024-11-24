package v1

import (
	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/masse/common"
)

type User common.User

func (user *User) SetCopyOption(info *llb.CopyInfo) {
	user.ChownOpt(info.ChownOpt).SetCopyOption(info)
}

func (user *User) SetMkfileOption(info *llb.MkfileInfo) {
	user.ChownOpt(info.ChownOpt).SetMkfileOption(info)
}

func (user *User) SetMkdirOption(info *llb.MkdirInfo) {
	user.ChownOpt(info.ChownOpt).SetMkdirOption(info)
}

func (user *User) UserOpt() *llb.UserOpt {
	opt := &llb.UserOpt{
		Name: user.User,
	}

	if user.UID != nil {
		opt.UID = int(*user.UID)
	}

	return opt
}

func (user *User) ChownOpt(other *llb.ChownOpt) *llb.ChownOpt {
	chown := &llb.ChownOpt{}

	if other != nil {
		chown.Group = other.Group
	}

	chown.User = user.UserOpt()

	return chown
}
