package v1

import (
	"cuelang.org/go/cue"
	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/masse/internal/lookup"
)

func (c *compiler) compileGit(state llb.State, v cue.Value) (llb.State, error) {
	repo, err := lookup.String(v, "git")
	if err != nil {
		return state, vError(v, err)
	}

	ref, err := lookup.String(v, "ref")
	if err != nil {
		return state, vError(v, err)
	}

	options, err := lookup.DecodeOptions[GitOptions](v)
	if err != nil {
		return state, vError(v, err)
	}

	return llb.Git(repo, ref, c.constraints(), options), nil
}

type GitOptions []*GitOption

type GitOption struct {
	*KeepGitDir
	*Constraint
}

func (opts GitOptions) SetGitOption(gi *llb.GitInfo) {
	for _, opt := range opts {
		opt.SetGitOption(gi)
	}
}

func (opt *GitOption) SetGitOption(gi *llb.GitInfo) {
	llbOpt, ok := oneof[llb.GitOption](opt)
	if ok {
		llbOpt.SetGitOption(gi)
	}
}

type KeepGitDir struct {
	KeepGitDir bool `json:"keepGitDir"`
}

func (kgd *KeepGitDir) SetGitOption(gi *llb.GitInfo) {
	gi.KeepGitDir = kgd.KeepGitDir
}
