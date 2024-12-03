package v1

import (
	"cuelang.org/go/cue"
	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/masse/internal/lookup"
)

func (c *compiler) compileRm(state llb.State, fa *llb.FileAction, v cue.Value) (*llb.FileAction, error) {
	path, err := lookup.String(v, "rm")
	if err != nil {
		return fa, vError(v, err)
	}

	options, err := lookup.DecodeOptions[RmOptions](v)
	if err != nil {
		return fa, vError(v, err)
	}

	return fa.Rm(path, options), nil
}

type RmOptions []*RmOption

func (opts RmOptions) SetRmOption(info *llb.RmInfo) {
	for _, opt := range opts {
		opt.SetRmOption(info)
	}
}

type RmOption struct {
	*AllowNotFound
	*AllowWildcard
}

func (opt *RmOption) SetRmOption(info *llb.RmInfo) {
	withOneOf(opt, func(ro llb.RmOption) { ro.SetRmOption(info) })
}

type AllowNotFound struct {
	AllowNotFound bool `json:"allowNotFound"`
}

func (opt *AllowNotFound) SetRmOption(info *llb.RmInfo) {
	llb.WithAllowNotFound(opt.AllowNotFound).SetRmOption(info)
}

type AllowWildcard struct {
	AllowWildcard bool `json:"allowWildcard"`
}

func (opt *AllowWildcard) SetRmOption(info *llb.RmInfo) {
	llb.WithAllowWildcard(opt.AllowWildcard).SetRmOption(info)
}
