package state

import "github.com/moby/buildkit/client/llb"

type Git struct {
	Repo    string     `json:"git"`
	Ref     string     `json:"ref"`
	Options GitOptions `json:"optionsValue"`
}

func (git *Git) CompileSource(_ ChainStates, constraints ...llb.ConstraintsOpt) (llb.State, error) {
	return llb.Git(
		git.Repo,
		git.Ref,
		append(constraintsTo[llb.GitOption](constraints), git.Options)...,
	), nil
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
