package state

import "github.com/moby/buildkit/client/llb"

type LLBConstraints interface {
	LLBConstraints() []llb.ConstraintsOpt
}

type Constraint struct {
	*Platform
}

func (c *Constraint) LLBConstraints() []llb.ConstraintsOpt {
	llbCons, ok := oneof[LLBConstraints](c)
	if ok {
		return llbCons.LLBConstraints()
	}

	return []llb.ConstraintsOpt{}
}

func (c *Constraint) SetImageOption(info *llb.ImageInfo) {
	for _, cons := range c.LLBConstraints() {
		cons.SetImageOption(info)
	}
}

func (c *Constraint) SetLocalOption(info *llb.LocalInfo) {
	for _, cons := range c.LLBConstraints() {
		cons.SetLocalOption(info)
	}
}

func (c *Constraint) SetGitOption(info *llb.GitInfo) {
	for _, cons := range c.LLBConstraints() {
		cons.SetGitOption(info)
	}
}
