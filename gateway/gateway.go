package gateway

import (
	"context"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/gateway/client"
	"github.com/pkg/errors"
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
	data, err := gw.readFromLayoutLocal(gw.options.LayoutFile, true)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read layout")
	}
}

func (gw *Gateway) readFromLayoutLocal(filepath string, required bool) ([]byte, error) {
	st := llb.Local(gw.options.LayoutLocal,
		llb.SessionID(gw.options.SessionID),
		llb.FollowPaths([]string{filepath}),
		llb.SharedKeyHint(gw.options.LayoutLocal+"-"+filepath),
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
