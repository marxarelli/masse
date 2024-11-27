package v1

import (
	"cuelang.org/go/cue"
	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/masse/internal/lookup"
)

func (c *compiler) compileImage(state llb.State, v cue.Value) (llb.State, error) {
	ref, err := lookup.String(v, "image")
	if err != nil {
		return state, vError(v, err)
	}

	inherit, err := lookup.Bool(v, "inherit")
	if err != nil {
		return state, vError(v, err)
	}

	options, err := lookup.DecodeListOrSingle[ImageOptions](v, "options")
	if err != nil {
		return state, vError(v, err)
	}

	imageOptions := []llb.ImageOption{c.constraints(), options}

	if inherit {
		imageOptions = append(imageOptions, llb.WithMetaResolver(c.config.ImageMetaResolver))
	}

	return llb.Image(ref, imageOptions...), nil
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
