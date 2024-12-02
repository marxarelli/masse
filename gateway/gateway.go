package gateway

import (
	"context"

	"github.com/moby/buildkit/frontend/dockerui"
	"github.com/moby/buildkit/frontend/gateway/client"
	dockerspec "github.com/moby/docker-image-spec/specs-go/v1"
	oci "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	"gitlab.wikimedia.org/dduvall/masse/common"
	v1compiler "gitlab.wikimedia.org/dduvall/masse/compiler/v1"
	"gitlab.wikimedia.org/dduvall/masse/config"
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

			// Solve LLB state to client result
			res, err := gw.bkClient.Solve(ctx, client.SolveRequest{
				Definition:   def.ToPB(),
				CacheImports: gw.CacheImports,
			})
			if err != nil {
				return nil, nil, nil, err
			}

			ref, err := res.SingleRef()
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

func (gw *Gateway) loadConfig() (*config.Root, error) {
	src, err := gw.ReadEntrypoint(gw.ctx, "CUE")
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config")
	}

	root, err := config.Load(src.Filename, src.Data, gw.Config.BuildArgs)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load config %q", src.Filename)
	}

	return root, nil
}
