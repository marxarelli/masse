package v1

import "github.com/moby/buildkit/client/llb"

type IgnoreCache struct {
	IgnoreCache bool `json:"ignoreCache"`
}

func (ic *IgnoreCache) constraintsOpt() llb.ConstraintsOpt {
	if ic.IgnoreCache {
		return llb.IgnoreCache
	}

	return NoopConstraint()
}

func (ic *IgnoreCache) SetConstraintsOption(c *llb.Constraints) {
	ic.constraintsOpt().SetConstraintsOption(c)
}

func (ic *IgnoreCache) SetRunOption(info *llb.ExecInfo) {
	ic.constraintsOpt().SetRunOption(info)
}

func (ic *IgnoreCache) SetImageOption(info *llb.ImageInfo) {
	ic.constraintsOpt().SetImageOption(info)
}

func (ic *IgnoreCache) SetLocalOption(info *llb.LocalInfo) {
	ic.constraintsOpt().SetLocalOption(info)
}

func (ic *IgnoreCache) SetGitOption(info *llb.GitInfo) {
	ic.constraintsOpt().SetGitOption(info)
}

func (ic *IgnoreCache) SetHTTPOption(info *llb.HTTPInfo) {
	ic.constraintsOpt().SetHTTPOption(info)
}

func (ic *IgnoreCache) SetOCILayoutOption(info *llb.OCILayoutInfo) {
	ic.constraintsOpt().SetOCILayoutOption(info)
}
