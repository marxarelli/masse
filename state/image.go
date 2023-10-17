package state

import (
	"github.com/moby/buildkit/client/llb"
)

type Image struct {
	Ref     string       `json:"image"`
	Options ImageOptions `json:"options"`
}

func (image *Image) CompileSource(_ ChainStates, constraints ...llb.ConstraintsOpt) (llb.State, error) {
	return llb.Image(
		image.Ref,
		append(constraintsTo[llb.ImageOption](constraints), image.Options)...,
	), nil
}

type ImageOptions []*ImageOption

func (opts ImageOptions) SetImageOption(info *llb.ImageInfo) {
	for _, opt := range opts {
		opt.SetImageOption(info)
	}
}

type ImageOption struct {
	*LayerLimit
	*Constraint
}

func (opt *ImageOption) SetImageOption(info *llb.ImageInfo) {
	llbOpt, ok := oneof[llb.ImageOption](opt)
	if ok {
		llbOpt.SetImageOption(info)
	}
}

type LayerLimit struct {
	LayerLimit uint32 `json:"layerLimit"`
}

func (ll *LayerLimit) SetImageOption(info *llb.ImageInfo) {
	llb.WithLayerLimit(int(ll.LayerLimit)).SetImageOption(info)
}
