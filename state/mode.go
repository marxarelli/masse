package state

import (
	"os"

	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/phyton/common"
)

type Mode common.Mode

func (mode *Mode) SetCopyOption(info *llb.CopyInfo) {
	info.Mode = mode.FileMode()
}

func (mode *Mode) FileMode() *os.FileMode {
	fm := os.FileMode(mode.Mode)
	return &fm
}
