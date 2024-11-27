package v1

import "github.com/moby/buildkit/client/llb"

type Constraints []*Constraint

func (cons Constraints) SetConstraintsOption(constraints *llb.Constraints) {
	for _, c := range cons {
		c.SetConstraintsOption(constraints)
	}
}

func (cons Constraints) SetRunOption(info *llb.ExecInfo) {
	for _, c := range cons {
		c.SetRunOption(info)
	}
}

func (cons Constraints) SetImageOption(info *llb.ImageInfo) {
	for _, c := range cons {
		c.SetImageOption(info)
	}
}

func (cons Constraints) SetLocalOption(info *llb.LocalInfo) {
	for _, c := range cons {
		c.SetLocalOption(info)
	}
}

func (cons Constraints) SetGitOption(info *llb.GitInfo) {
	for _, c := range cons {
		c.SetGitOption(info)
	}
}

func (cons Constraints) SetHTTPOption(info *llb.HTTPInfo) {
	for _, c := range cons {
		c.SetHTTPOption(info)
	}
}

func (cons Constraints) SetOCILayoutOption(info *llb.OCILayoutInfo) {
	for _, c := range cons {
		c.SetOCILayoutOption(info)
	}
}

type Constraint struct {
	*Platform
}

func (c *Constraint) SetConstraintsOption(constraints *llb.Constraints) {
	withOneOf(c, func(opt llb.ConstraintsOpt) { opt.SetConstraintsOption(constraints) })
}

func (c *Constraint) SetRunOption(info *llb.ExecInfo) {
	withOneOf(c, func(opt llb.RunOption) { opt.SetRunOption(info) })
}

func (c *Constraint) SetImageOption(info *llb.ImageInfo) {
	withOneOf(c, func(opt llb.ImageOption) { opt.SetImageOption(info) })
}

func (c *Constraint) SetLocalOption(info *llb.LocalInfo) {
	withOneOf(c, func(opt llb.LocalOption) { opt.SetLocalOption(info) })
}

func (c *Constraint) SetGitOption(info *llb.GitInfo) {
	withOneOf(c, func(opt llb.GitOption) { opt.SetGitOption(info) })
}

func (c *Constraint) SetHTTPOption(info *llb.HTTPInfo) {
	withOneOf(c, func(opt llb.HTTPOption) { opt.SetHTTPOption(info) })
}

func (c *Constraint) SetOCILayoutOption(info *llb.OCILayoutInfo) {
	withOneOf(c, func(opt llb.OCILayoutOption) { opt.SetOCILayoutOption(info) })
}
