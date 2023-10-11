package state

import (
	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/phyton/common"
)

type Host common.Host

func (host *Host) LLBRunOptions(_ ChainStates) ([]llb.RunOption, error) {
	return []llb.RunOption{
		llb.AddExtraHost(host.Host, host.IP),
	}, nil
}
