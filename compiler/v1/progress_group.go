package v1

import "github.com/moby/buildkit/client/llb"

type ProgressGroup struct {
	ProgressGroup string `json:"progressGroup"`
	Name          string `json:"name"`
	Weak          bool   `json:"weak"`
}

func (pg *ProgressGroup) constraintsOpt() llb.ConstraintsOpt {
	return llb.ProgressGroup(pg.ProgressGroup, pg.Name, pg.Weak)
}

func (pg *ProgressGroup) SetConstraintsOption(c *llb.Constraints) {
	pg.constraintsOpt().SetConstraintsOption(c)
}

func (pg *ProgressGroup) SetRunOption(info *llb.ExecInfo) {
	pg.constraintsOpt().SetRunOption(info)
}

func (pg *ProgressGroup) SetImageOption(info *llb.ImageInfo) {
	pg.constraintsOpt().SetImageOption(info)
}

func (pg *ProgressGroup) SetLocalOption(info *llb.LocalInfo) {
	pg.constraintsOpt().SetLocalOption(info)
}

func (pg *ProgressGroup) SetGitOption(info *llb.GitInfo) {
	pg.constraintsOpt().SetGitOption(info)
}

func (pg *ProgressGroup) SetHTTPOption(info *llb.HTTPInfo) {
	pg.constraintsOpt().SetHTTPOption(info)
}

func (pg *ProgressGroup) SetOCILayoutOption(info *llb.OCILayoutInfo) {
	pg.constraintsOpt().SetOCILayoutOption(info)
}
