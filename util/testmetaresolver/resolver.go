package testmetaresolver

import (
	"context"
	"encoding/json"

	"github.com/moby/buildkit/client/llb"
	digest "github.com/opencontainers/go-digest"
	oci "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

type testResolver struct {
	ref    string
	digest digest.Digest
	image  oci.Image
}

// New returns a noop [llb.ImageMetaResolver] implementation that resolves any
// request for the given image ref to a the given [oci.Image].
func New(ref string, image oci.Image) llb.ImageMetaResolver {
	return &testResolver{
		ref:    ref,
		digest: digest.FromBytes([]byte(ref)),
		image:  image,
	}
}

// ResolveImageConfig returns the [digest.Digest] and [oci.Image] (marshaled
// to JSON) given at creation.
func (tr *testResolver) ResolveImageConfig(ctx context.Context, ref string, opt llb.ResolveImageConfigOpt) (string, digest.Digest, []byte, error) {
	image := oci.Image{
		Created: tr.image.Created,
		Author:  tr.image.Author,
		Platform: oci.Platform{
			Architecture: tr.image.Architecture,
			OS:           tr.image.OS,
		},
		Config:  tr.image.Config,
		RootFS:  tr.image.RootFS,
		History: tr.image.History,
	}

	if opt.Platform != nil {
		image.Architecture = opt.Platform.Architecture
		image.OS = opt.Platform.OS
		image.Variant = opt.Platform.Variant
	}

	cfg, err := json.Marshal(image)
	if err != nil {
		return "", "", nil, errors.WithStack(err)
	}

	return tr.ref, tr.digest, cfg, nil
}
