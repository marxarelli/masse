package v1

import "github.com/moby/buildkit/client/llb"

type CustomName struct {
	CustomName string `json:"customName"`
}

func (cn *CustomName) constraintsOpt() llb.ConstraintsOpt {
	return llb.WithCustomName(cn.CustomName)
}

func (cn *CustomName) SetConstraintsOption(c *llb.Constraints) {
	cn.constraintsOpt().SetConstraintsOption(c)
}

func (cn *CustomName) SetRunOption(info *llb.ExecInfo) {
	cn.constraintsOpt().SetRunOption(info)
}

func (cn *CustomName) SetImageOption(info *llb.ImageInfo) {
	cn.constraintsOpt().SetImageOption(info)
}

func (cn *CustomName) SetLocalOption(info *llb.LocalInfo) {
	cn.constraintsOpt().SetLocalOption(info)
}

func (cn *CustomName) SetGitOption(info *llb.GitInfo) {
	cn.constraintsOpt().SetGitOption(info)
}

func (cn *CustomName) SetHTTPOption(info *llb.HTTPInfo) {
	cn.constraintsOpt().SetHTTPOption(info)
}

func (cn *CustomName) SetOCILayoutOption(info *llb.OCILayoutInfo) {
	cn.constraintsOpt().SetOCILayoutOption(info)
}
