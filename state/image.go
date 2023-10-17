package state

import (
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/client/llb/imagemetaresolver"
)

type Image struct {
	Ref     string       `json:"image"`
	Inherit bool         `json:"inherit"`
	Options ImageOptions `json:"options"`
}

func (image *Image) CompileSource(_ ChainStates, constraints ...llb.ConstraintsOpt) (llb.State, error) {
	options := append(constraintsTo[llb.ImageOption](constraints), image.Options)

	if image.Inherit {
		options = append(
			options,
			llb.WithMetaResolver(imagemetaresolver.Default()),
		)
	}

	return llb.Image(
		image.Ref,
		options...,
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
