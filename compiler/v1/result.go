package v1

import (
	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/masse/common"
)

type result struct {
	ref      string
	platform common.Platform
	state    llb.State
	deps     map[string]llb.State
}

func (res *result) ChainRef() string {
	return res.ref
}

func (res *result) ChainState() llb.State {
	return res.state
}

func (res *result) Platform() common.Platform {
	return res.platform
}

func (res *result) DependencyChainStates() map[string]llb.State {
	return res.deps
}
