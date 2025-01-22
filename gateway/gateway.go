package gateway

import (
	"context"
	"encoding/json"
	"path"
	"strings"

	"cuelang.org/go/mod/modconfig"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/dockerui"
	"github.com/moby/buildkit/frontend/gateway/client"
	dockerspec "github.com/moby/docker-image-spec/specs-go/v1"
	oci "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	"gitlab.wikimedia.org/dduvall/masse/common"
	v1compiler "gitlab.wikimedia.org/dduvall/masse/compiler/v1"
	"gitlab.wikimedia.org/dduvall/masse/config"
	"gitlab.wikimedia.org/dduvall/masse/load"
)

func Run(ctx context.Context, c client.Client) (*client.Result, error) {
	gw, err := New(ctx, c)
	if err != nil {
		return nil, err
	}

	return gw.Run()
}

// Gateway reads in a client layout and transforms the target into LLB for
// buildkitd
type Gateway struct {
	*dockerui.Client
	bkClient client.Client
	ctx      context.Context
}

func New(ctx context.Context, bkClient client.Client) (*Gateway, error) {
	client, err := dockerui.NewClient(bkClient)
	if err != nil {
		return nil, err
	}

	return &Gateway{
		Client:   client,
		bkClient: bkClient,
		ctx:      ctx,
	}, nil
}

func (gw *Gateway) Run() (*client.Result, error) {
	root, err := gw.loadConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load config")
	}

	target, ok := root.Targets[gw.Config.Target]
	if !ok {
		return nil, errors.Wrapf(err, "unknown target %q", gw.Config.Target)
	}

	if gw.Config.TargetPlatforms == nil || len(gw.Config.TargetPlatforms) == 0 {
		gw.Config.TargetPlatforms = target.OCIPlatforms()
	} else {
		// TODO validate requested platforms against target
	}

	resultBuilder, err := gw.Build(
		gw.ctx,
		func(ctx context.Context, platform *oci.Platform, idx int) (
			client.Reference,
			*dockerspec.DockerOCIImage,
			*dockerspec.DockerOCIImage,
			error,
		) {
			targetPlatform := common.PlatformFromOCI(platform)

			compiler := v1compiler.New(
				root.Chains,
				v1compiler.WithPlatform(targetPlatform),
				v1compiler.WithContext(ctx),
			)

			// Compile to LLB state
			st, err := compiler.Compile(target)
			if err != nil {
				return nil, nil, nil, errors.Wrap(err, "failed to compile target")
			}

			def, err := st.Marshal(ctx)
			if err != nil {
				return nil, nil, nil, errors.Wrap(err, "failed to marshal target state")
			}

			ref, err := gw.solve(client.SolveRequest{
				Definition:   def.ToPB(),
				CacheImports: gw.CacheImports,
			})

			if err != nil {
				return nil, nil, nil, err
			}

			img := target.NewImage(targetPlatform)

			dimg := dockerspec.DockerOCIImage{
				Image: img,
				Config: dockerspec.DockerOCIImageConfig{
					ImageConfig: img.Config,
				},
			}

			return ref, &dimg, nil, nil
		},
	)

	if err != nil {
		return nil, err
	}

	return resultBuilder.Finalize()
}

func (gw *Gateway) parameters() map[string]string {
	params := map[string]string{}

	for key, value := range gw.Config.BuildArgs {

		if trimmed, found := strings.CutPrefix(key, "PARAMETER_"); found {
			params[trimmed] = value
		}
	}

	return params
}

func (gw *Gateway) solve(req client.SolveRequest) (client.Reference, error) {
	var nilRef client.Reference

	res, err := gw.bkClient.Solve(gw.ctx, req)
	if err != nil {
		return nilRef, err
	}

	return res.SingleRef()
}

// loadConfig solves the config source llb.State and reads the CUE config file
// as well as any cue.mod/module.cue file relative to the directory of the
// config file path.
func (gw *Gateway) loadConfig() (*config.Root, error) {
	// XXX until we can read multiple entrypoint files, we'll have to hack the
	// follow paths and perform a second solve below. If
	// https://github.com/moby/buildkit/pull/5866 is merged this hack can be
	// removed in favor of ReadEntrypointFiles.
	src, err := gw.ReadEntrypoint(gw.ctx, "CUE", followExtraPath("../cue.mod/module.cue"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config")
	}

	// perform second solve to retrieve cue.mod/module.cue
	ref, err := gw.solve(client.SolveRequest{
		Definition: src.Definition.ToPB(),
	})

	moduleOverlay := map[string][]byte{}

	if err == nil {
		moduleFile := path.Join(src.Filename, "../cue.mod/module.cue")

		dt, err := ref.ReadFile(gw.ctx, client.ReadRequest{Filename: moduleFile})
		if err == nil {
			moduleOverlay[moduleFile] = dt
		}
	}

	env := []string{}
	for key, value := range gw.Config.BuildArgs {
		env = append(env, key+"="+value)
	}

	registry, err := modconfig.NewRegistry(&modconfig.Config{
		Transport: NewRegistryTransport(gw),
		Env:       env,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create CUE registry")
	}

	root, err := config.Load(
		src.Filename,
		src.Data,
		gw.parameters(),
		load.WithRegistry(registry),
		load.WithOverlayFiles(moduleOverlay),
		load.WithEnvMap(gw.Config.BuildArgs),
	)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to load config %q", src.Filename)
	}

	return root, nil
}

type localOptionFunc func(*llb.LocalInfo)

func (lo localOptionFunc) SetLocalOption(info *llb.LocalInfo) {
	lo(info)
}

func followExtraPath(xpath string) llb.LocalOption {
	return localOptionFunc(func(info *llb.LocalInfo) {
		var paths []string
		if json.Unmarshal([]byte(info.FollowPaths), &paths) == nil && len(paths) > 0 {
			filename := paths[0]
			paths = append(paths, path.Join(filename, xpath))
			dt, err := json.Marshal(paths)
			if err == nil {
				info.FollowPaths = string(dt)
			}
		}
	})
}
