package state

import (
	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/phyton/common"
)

type User common.User

func (user *User) SetCopyOption(info *llb.CopyInfo) {
	chown := &llb.ChownOpt{}
	if info.ChownOpt == nil {
		chown = info.ChownOpt
	}

	chown.User = user.UserOpt()
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
