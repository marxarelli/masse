package state

import (
	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/phyton/common"
)

type Platform struct {
	Platform common.Platform `json:"platformValue"`
}

func (p *Platform) LLBConstraints() []llb.ConstraintsOpt {
	return []llb.ConstraintsOpt{
		llb.Platform(p.Platform.OCI()),
	}
}
