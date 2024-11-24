package v1

import (
	"os"

	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/masse/common"
)

const (
	defaultFileMode uint32 = uint32(0o0644)
	defaultDirMode  uint32 = uint32(0o0755)
)

type Mode common.Mode

func (mode *Mode) SetCopyOption(info *llb.CopyInfo) {
	mode.ChmodOpt().SetCopyOption(info)
}

func (mode *Mode) FileMode() os.FileMode {
	return os.FileMode(mode.Mode)
}

func (mode *Mode) ChmodOpt() llb.ChmodOpt {
	return llb.ChmodOpt{
		Mode: mode.FileMode(),
	}
}

func DefaultFileMode() *Mode {
	return &Mode{Mode: defaultFileMode}
}

func DefaultDirMode() *Mode {
	return &Mode{Mode: defaultDirMode}
}
