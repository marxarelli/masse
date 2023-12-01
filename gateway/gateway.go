package gateway

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/exporter/containerimage/exptypes"
	"github.com/moby/buildkit/frontend/gateway/client"
	"github.com/pkg/errors"
	"gitlab.wikimedia.org/dduvall/phyton/common"
	"gitlab.wikimedia.org/dduvall/phyton/config"
	"gitlab.wikimedia.org/dduvall/phyton/state"
	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context, c client.Client) (*client.Result, error) {
	opts, err := ParseOptions(c.BuildOpts())

	if err != nil {
		return nil, errors.Wrap(err, "failed to parse gateway options")
	}

	return New(ctx, c, opts).Run()
}

// Gateway reads in a client layout and transforms the target into LLB for
// buildkitd
type Gateway struct {
	client.Client
	ctx     context.Context
	options *Options
}

func New(ctx context.Context, c client.Client, opts *Options) *Gateway {
	return &Gateway{Client: c, ctx: ctx, options: opts}
}

func (gw *Gateway) Run() (*client.Result, error) {
	data, err := gw.readFromConfigLocal(gw.options.ConfigFile, true)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read config from %s", gw.options.ConfigFile)
	}

	root, err := config.Load(gw.options.ConfigFile, data)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load config from %s", gw.options.ConfigFile)
	}

	target, ok := root.Targets[gw.options.Target]
	if !ok {
		return nil, errors.Wrapf(err, "unknown target %q", gw.options.Target)
	}

	graph, err := target.Graph(root.Chains)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get target graph")
	}

	exportPlatforms := &exptypes.Platforms{
		Platforms: make([]exptypes.Platform, len(target.Platforms)),
	}
	finalResult := client.NewResult()

	eg, ctx := errgroup.WithContext(gw.ctx)

	// Solve for all platforms in parallel
	for i, tp := range target.Platforms {
		func(i int, tp common.Platform) {
			eg.Go(func() (err error) {
				solver := state.NewPlatformSolver(tp)

				// Solve to LLB state
				st, err := solver.Solve(graph)
				if err != nil {
					return errors.Wrap(err, "failed to solve target graph")
				}

				def, err := st.Marshal(ctx)
				if err != nil {
					return errors.Wrap(err, "failed to marshal LLB state")
				}

				// Solve LLB state to client result
				res, err := gw.Solve(ctx, client.SolveRequest{
					Definition:   def.ToPB(),
					CacheImports: gw.options.CacheOptions,
				})
				if err != nil {
					return err
				}

				ref, err := res.SingleRef()
				if err != nil {
					return err
				}

				imageJSON, err := json.Marshal(target.NewImage(tp))
				if err != nil {
					return errors.Wrap(err, "failed to marshal image config")
				}

				exportPlatform := tp.Export()
				exportPlatforms.Platforms[i] = exportPlatform

				finalResult.AddMeta(
					fmt.Sprintf("%s/%s", exptypes.ExporterImageConfigKey, exportPlatform),
					imageJSON,
				)
				finalResult.AddRef(exportPlatform.ID, ref)

				return nil
			})
		}(i, tp)
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	platformData, err := json.Marshal(exportPlatforms)
	if err != nil {
		return nil, err
	}

	finalResult.AddMeta(exptypes.ExporterPlatformsKey, platformData)

	return finalResult, nil
}

func (gw *Gateway) readFromConfigLocal(filepath string, required bool) ([]byte, error) {
	st := llb.Local(gw.options.ConfigLocal,
		llb.SessionID(gw.options.SessionID),
		llb.FollowPaths([]string{filepath}),
		llb.SharedKeyHint(gw.options.ConfigLocal+"-"+filepath),
	)

	def, err := st.Marshal(gw.ctx)
	if err != nil {
		return nil, err
	}

	res, err := gw.Solve(gw.ctx, client.SolveRequest{
		Definition: def.ToPB(),
	})
	if err != nil {
		return nil, err
	}

	ref, err := res.SingleRef()
	if err != nil {
		return nil, err
	}

	// If the file is not required, try to stat it first, and if it doesn't
	// exist, simply return an empty byte slice. If the file is required, we'll
	// save an extra stat call and just try to read it.
	if !required {
		_, err := ref.StatFile(gw.ctx, client.StatRequest{
			Path: filepath,
		})

		if err != nil {
			return []byte{}, nil
		}
	}

	fileBytes, err := ref.ReadFile(gw.ctx, client.ReadRequest{
		Filename: filepath,
	})

	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}
