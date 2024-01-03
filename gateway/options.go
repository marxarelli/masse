package gateway

import (
	"encoding/json"
	"strings"

	"github.com/containerd/containerd/platforms"
	controlapi "github.com/moby/buildkit/api/services/control"
	"github.com/moby/buildkit/frontend/gateway/client"
	oci "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"

	"gitlab.wikimedia.org/dduvall/masse/config"
)

const (
	keyCacheFrom    = "cache-from"    // for registry only. deprecated in favor of keyCacheImports
	keyCacheImports = "cache-imports" // JSON representation of []CacheOptionsEntry

	keyConfigLocal = "config-local"
	keyConfigFile  = "config"
	keyTarget      = "target"

	defaultConfigLocal = "config"
	defaultConfigFile  = "masse.cue"

	// Dockerfile syntax= compatibility
	dockerfileLocal   = "dockerfile" // tried prior to defaultLayoutLocal
	keyDockerfilePath = "filename"   // = keyConfigFile

	// Support the same build-arg: option prefix that buildkit's dockerfile
	// frontend supports. Use the values as Masse parameters.
	buildArgPrefix  = "build-arg:"
	parameterPrefix = "parameter:"
)

// Options stores options to configure the build process.
type Options struct {
	// Name of the client's local context that contains the config file
	ConfigLocal string

	// Path to the config file, relative to the ConfigLocal
	ConfigFile string

	// Target is the targets entry to build
	Target string

	// Parameters are user supplied build parameters
	Parameters config.Parameters

	// Session ID
	SessionID string

	// CacheOptions specifies caches to be imported prior to the build
	CacheOptions []client.CacheOptionsEntry
}

// NewOptions creates a new Options with default values assigned
func NewOptions() *Options {
	return &Options{
		ConfigLocal:  defaultConfigLocal,
		ConfigFile:   defaultConfigFile,
		Parameters:   config.Parameters{},
		CacheOptions: []client.CacheOptionsEntry{},
	}
}

// ParseOptions parses and returns a newly created Options from the given
// client options.
func ParseOptions(cbo client.BuildOpts) (*Options, error) {
	opts := NewOptions()

	var err error

	// Assume Dockerfile syntax= usage based on product
	// TODO test this
	if cbo.Product == "docker" {
		opts.ConfigLocal = dockerfileLocal
	}

	for k, v := range cbo.Opts {
		switch k {
		case keyConfigLocal:
			opts.ConfigLocal = v

		case keyTarget:
			opts.Target = v

		case keyConfigFile, keyDockerfilePath:
			opts.ConfigFile = v
		}
	}

	opts.CacheOptions, err = parseCacheOptions(cbo.Opts)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse cache options")
	}

	opts.SessionID = cbo.SessionID

	opts.Parameters = filterOpts(cbo.Opts, buildArgPrefix)
	for key, val := range filterOpts(cbo.Opts, parameterPrefix) {
		opts.Parameters[key] = val
	}

	return opts, nil
}

// parseCacheOptions handles given cache imports. Note that clients may give
// these options in two different ways, either as `cache-imports` or
// `cache-from`. The latter is used for registry based cache imports.
// See https://github.com/moby/buildkit/blob/v0.10/client/solve.go#L477
//
// TODO the master branch of buildkit removes the legacy `cache-from` key, so
// we should be able to eventually deprecate it, but that will involve
// dropping support for older buildctl and docker buildx clients.
func parseCacheOptions(opts map[string]string) ([]client.CacheOptionsEntry, error) {
	var cacheImports []client.CacheOptionsEntry
	// new API
	if cacheImportsStr := opts[keyCacheImports]; cacheImportsStr != "" {
		var cacheImportsUM []controlapi.CacheOptionsEntry
		if err := json.Unmarshal([]byte(cacheImportsStr), &cacheImportsUM); err != nil {
			return nil, errors.Wrapf(err, "failed to unmarshal %s (%q)", keyCacheImports, cacheImportsStr)
		}
		for _, um := range cacheImportsUM {
			cacheImports = append(cacheImports, client.CacheOptionsEntry{Type: um.Type, Attrs: um.Attrs})
		}
	}
	// old API
	if cacheFromStr := opts[keyCacheFrom]; cacheFromStr != "" {
		cacheFrom := strings.Split(cacheFromStr, ",")
		for _, s := range cacheFrom {
			im := client.CacheOptionsEntry{
				Type: "registry",
				Attrs: map[string]string{
					"ref": s,
				},
			}
			// FIXME(AkihiroSuda): skip append if already exists
			cacheImports = append(cacheImports, im)
		}
	}

	return cacheImports, nil
}

func parsePlatforms(v string) ([]*oci.Platform, error) {
	var pp []*oci.Platform
	for _, v := range strings.Split(v, ",") {
		p, err := platforms.Parse(v)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse target platform %s", v)
		}
		p = platforms.Normalize(p)
		pp = append(pp, &p)
	}
	return pp, nil
}

func filterOpts(opts map[string]string, prefix string) map[string]string {
	filtered := map[string]string{}

	for k, v := range opts {
		if strings.HasPrefix(k, prefix) {
			filtered[strings.TrimPrefix(k, prefix)] = v
		}
	}

	return filtered
}
