package v1

import (
	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/masse/common"
)

type Platform struct {
	Platform common.Platform `json:"platformValue"`
}

func (p *Platform) LLBConstraints() []llb.ConstraintsOpt {
	return []llb.ConstraintsOpt{
		llb.Platform(p.Platform.OCI()),
	}
}
