package target

import (
	"time"

	"cuelang.org/go/cue"
	oci "github.com/opencontainers/image-spec/specs-go/v1"
	"gitlab.wikimedia.org/dduvall/masse/common"
)

type Targets map[string]*Target

type Target struct {
	Build     cue.Value
	Platforms []common.Platform `json:"platformsValue"`
	Runtime   Runtime           `json:"runtime"`
	Labels    map[string]string `json:"labels"`
}

// UnmarshalCUE parses the given cue.Value into the target.
func (target *Target) UnmarshalCUE(v cue.Value) error {
	return v.Decode(target)
}

// NewImage returns the target as a new [oci.Image] for the given platform.
func (target *Target) NewImage(platform common.Platform) oci.Image {
	now := time.Now()

	return oci.Image{
		Created:  &now,
		Platform: platform.OCI(),
		Config:   target.Runtime.ImageConfig(),
	}
}

func (target *Target) OCIPlatforms() []oci.Platform {
	platforms := make([]oci.Platform, len(target.Platforms))
	for i, tp := range target.Platforms {
		platforms[i] = tp.OCI()
	}

	return platforms
}
