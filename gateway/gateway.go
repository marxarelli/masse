package gateway

import (
	"context"
	"encoding/json"
	"path"
	"strings"
	"sync"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/client/llb/sourceresolver"
	"github.com/moby/buildkit/frontend"
	"github.com/moby/buildkit/frontend/attestations/sbom"
	"github.com/moby/buildkit/frontend/dockerui"
	"github.com/moby/buildkit/frontend/gateway/client"
	"github.com/moby/buildkit/solver/pb"
	"github.com/moby/buildkit/solver/result"
	dockerspec "github.com/moby/docker-image-spec/specs-go/v1"
	oci "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	"gitlab.wikimedia.org/dduvall/masse/common"
	v1compiler "gitlab.wikimedia.org/dduvall/masse/compiler/v1"
	"gitlab.wikimedia.org/dduvall/masse/config"
	"gitlab.wikimedia.org/dduvall/masse/load"
	"gitlab.wikimedia.org/dduvall/masse/target"
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

	buildTarget, ok := root.Targets[gw.Config.Target]
	if !ok {
		return nil, errors.Wrapf(err, "unknown target %q", gw.Config.Target)
	}

	if gw.Config.TargetPlatforms == nil || len(gw.Config.TargetPlatforms) == 0 {
		gw.Config.TargetPlatforms = buildTarget.OCIPlatforms()
	} else {
		// TODO validate requested platforms against target
	}

	var scanner sbom.Scanner

	if gw.SBOM != nil {
		scanner, err = gw.newSBOMScanner(buildTarget)
		if err != nil {
			return nil, err
		}
	}

	scanTargets := sync.Map{}

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
				v1compiler.WithIgnoreCache(gw.ignoreCache()),
				v1compiler.WithMainContextLoader(func(ctx context.Context, opts ...llb.LocalOption) (*llb.State, error) {
					return gw.MainContext(ctx, opts...)
				}),
				v1compiler.WithNamedContextLoader(func(ctx context.Context, name string, opts ...llb.LocalOption) (*llb.State, error) {
					contextOpt := dockerui.ContextOpt{
						AsyncLocalOpts: func() []llb.LocalOption { return opts },
						Platform:       platform,
						ResolveMode:    resolveModeName(gw.ImageResolveMode),
					}
					namedContext, err := gw.NamedContext(name, contextOpt)
					if err != nil {
						return nil, err
					}

					state, _, err := namedContext.Load(ctx)
					return state, err
				}),
			)

			// Compile to LLB state
			compileResult, err := compiler.Compile(buildTarget)
			if err != nil {
				return nil, nil, nil, errors.Wrap(err, "failed to compile target")
			}

			def, err := compileResult.ChainState().Marshal(ctx)
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

			img := buildTarget.NewImage(targetPlatform)

			dimg := dockerspec.DockerOCIImage{
				Image: img,
				Config: dockerspec.DockerOCIImageConfig{
					ImageConfig: img.Config,
				},
			}

			scanTargets.Store(compileResult.Platform().ID(), compileResult)

			return ref, &dimg, nil, nil
		},
	)

	if err != nil {
		return nil, err
	}

	if scanner != nil {
		err = resultBuilder.EachPlatform(gw.ctx, func(ctx context.Context, id string, _ oci.Platform) error {
			v, ok := scanTargets.Load(id)
			if !ok {
				return errors.Errorf("no scan targets for %s", id)
			}

			compileResult, ok := v.(target.CompilerResult)
			if !ok {
				return errors.Errorf("invalid scan targets for %T", v)
			}

			att, err := scanner(
				ctx,
				id,
				compileResult.ChainState(),
				compileResult.DependencyChainStates(),
			)
			if err != nil {
				return err
			}

			attSolve, err := result.ConvertAttestation(&att, func(st *llb.State) (client.Reference, error) {
				def, err := st.Marshal(ctx)
				if err != nil {
					return nil, err
				}
				r, err := gw.bkClient.Solve(ctx, frontend.SolveRequest{
					Definition: def.ToPB(),
				})
				if err != nil {
					return nil, err
				}
				return r.Ref, nil
			})
			if err != nil {
				return err
			}
			resultBuilder.AddAttestation(id, *attSolve)
			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	return resultBuilder.Finalize()
}

func (gw *Gateway) newSBOMScanner(target *target.Target) (sbom.Scanner, error) {
	targetSBOM := target.Attestations.SBOM

	generator := gw.SBOM.Generator
	if targetSBOM.Generator != "" {
		generator = targetSBOM.Generator
	}

	parameters := map[string]string{}
	if targetSBOM.Parameters != nil {
		for name, value := range targetSBOM.Parameters {
			parameters[name] = value
		}
	}
	if gw.SBOM.Parameters != nil {
		for name, value := range gw.SBOM.Parameters {
			parameters[name] = value
		}
	}

	scanner, err := sbom.CreateSBOMScanner(
		gw.ctx, gw.bkClient, generator,
		sourceresolver.Opt{
			ImageOpt: &sourceresolver.ResolveImageOpt{
				ResolveMode: resolveModeName(gw.ImageResolveMode),
			},
		},
		parameters,
	)

	if err != nil {
		return scanner, err
	}

	return func(
		ctx context.Context,
		name string,
		state llb.State,
		deps map[string]llb.State,
		opts ...llb.ConstraintsOpt,
	) (result.Attestation[*llb.State], error) {
		filteredDeps := map[string]llb.State{}
		scanRefs := targetSBOM.ScanChainRefMap()

		for ref, depState := range deps {
			if targetSBOM.ScanAll() || scanRefs[ref] {
				filteredDeps[ref] = depState
			}
		}

		return scanner(ctx, name, state, filteredDeps, opts...)
	}, nil
}

func (gw *Gateway) ignoreCache() bool {
	return gw.Client.IsNoCache(gw.Config.Target)
}

func (gw *Gateway) parameters() map[string]string {
	params := map[string]string{}

	for key, value := range gw.Config.BuildArgs {

		if trimmed, found := strings.CutPrefix(key, "masse:parameter:"); found {
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

	root, err := config.Load(
		src.Filename,
		src.Data,
		gw.parameters(),
		load.WithOverlayFiles(moduleOverlay),
		load.WithEnvMap(gw.Config.BuildArgs),
		load.WithRegistryTransport(NewRegistryTransport(gw)),
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

func resolveModeName(mode llb.ResolveMode) string {
	switch mode {
	case llb.ResolveModeForcePull:
		return pb.AttrImageResolveModeForcePull
	case llb.ResolveModePreferLocal:
		return pb.AttrImageResolveModePreferLocal
	}
	return pb.AttrImageResolveModeDefault
}
