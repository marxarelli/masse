package state

import (
	"github.com/moby/buildkit/client/llb"
	oci "github.com/opencontainers/image-spec/specs-go/v1"
	"gitlab.wikimedia.org/dduvall/phyton/common"
)

type Image struct {
	Ref     string       `json:"image"`
	Options ImageOptions `json:"options"`
}

func (image *Image) Compile(_ llb.State, _ ChainStates) (llb.State, error) {
	return llb.Image(image.Ref, image.Options), nil
}

type ImageOptions []*ImageOption

func (opts ImageOptions) SetImageOption(info *llb.ImageInfo) {
	for _, opt := range opts {
		opt.SetImageOption(info)
	}
}

type ImageOption struct {
	*Platform
	*LayerLimit
}

func (opt *ImageOption) SetImageOption(info *llb.ImageInfo) {
	llbOpt, ok := oneof[llb.ImageOption](opt)
	if ok {
		llbOpt.SetImageOption(info)
	}
}

type Platform struct {
	Platform common.Platform `json:"value"`
}

func (p *Platform) SetImageOption(info *llb.ImageInfo) {
	llb.Platform(oci.Platform{
		OS:           p.Platform.OS,
		Architecture: p.Platform.Architecture,
		Variant:      p.Platform.Variant,
	}).SetImageOption(info)
}

type LayerLimit struct {
	LayerLimit uint32 `json:"layerLimit"`
}

func (ll *LayerLimit) SetImageOption(info *llb.ImageInfo) {
	llb.WithLayerLimit(int(ll.LayerLimit)).SetImageOption(info)
}
