package v1

import (
	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/masse/common"
)

var defaultPlatform = &Platform{
	Platform: common.Platform{
		OS:           "linux",
		Architecture: "amd64",
	},
}

type Platform struct {
	Platform common.Platform `json:"platformValue"`
}

func (p *Platform) constraintsOpt() llb.ConstraintsOpt {
	return llb.Platform(p.Platform.OCI())
}

func (p *Platform) SetImageOption(info *llb.ImageInfo) {
	p.constraintsOpt().SetImageOption(info)
}

func (p *Platform) SetLocalOption(info *llb.LocalInfo) {
	p.constraintsOpt().SetLocalOption(info)
}

func (p *Platform) SetGitOption(info *llb.GitInfo) {
	p.constraintsOpt().SetGitOption(info)
}
