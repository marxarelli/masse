package common

import (
	"github.com/containerd/containerd/platforms"
	"github.com/moby/buildkit/exporter/containerimage/exptypes"
	oci "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

type Platform struct {
	OS           string `json:"os"`
	Architecture string `json:"architecture"`
	Variant      string `json:"variant"`
}

// ParsePlatform returns a new [Platform] for the given platform shorthand
// name (e.g. "linux/arm64/v8").
func ParsePlatform(platformName string) (Platform, error) {
	var platform Platform

	p, err := platforms.Parse(platformName)
	if err != nil {
		return platform, errors.Wrapf(err, "failed to parse platform %s", platformName)
	}

	p = platforms.Normalize(p)

	platform.OS = p.OS
	platform.Architecture = p.Architecture
	platform.Variant = p.Variant

	return platform, nil
}

func (p Platform) Export() exptypes.Platform {
	return exptypes.Platform{
		ID:       p.ID(),
		Platform: p.OCI(),
	}
}

func (p Platform) OCI() oci.Platform {
	return oci.Platform{
		OS:           p.OS,
		Architecture: p.Architecture,
		Variant:      p.Variant,
	}
}

func (p Platform) ID() string {
	return platforms.Format(p.OCI())
}

func (p Platform) String() string {
	return p.ID()
}
