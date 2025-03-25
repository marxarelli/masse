package gateway

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/gateway/client"
	"github.com/pkg/errors"
)

const (
	authSecretNameOptionPrefix = "CUE_REGISTRY_AUTH_SECRET."
	responseFilename           = "/http.response"
	progressGroupID            = "masse.modules"
	progressGroupName          = "⚪ fetch CUE modules"
)

type registryTransport struct {
	ctx       context.Context
	client    client.Client
	buildArgs map[string]string
	options   []llb.HTTPOption
}

func NewRegistryTransport(gw *Gateway) http.RoundTripper {
	opts := []llb.HTTPOption{
		llb.ProgressGroup(progressGroupID, progressGroupName, false),
	}

	if gw.ignoreCache() {
		opts = append(opts, llb.IgnoreCache)
	}

	return &registryTransport{
		ctx:       gw.ctx,
		client:    gw.bkClient,
		buildArgs: gw.Config.BuildArgs,
		options:   opts,
	}
}

func (rt *registryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Method != "GET" {
		return nil, errors.New("can only perform GET requests")
	}

	def, err := rt.requestToLLB(req).Marshal(rt.ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal HTTP LLB state")
	}

	res, err := rt.client.Solve(rt.ctx, client.SolveRequest{
		Definition: def.ToPB(),
	})
	if err != nil {
		return nil, err
	}

	ref, err := res.SingleRef()
	if err != nil {
		return nil, err
	}

	dt, err := ref.ReadFile(rt.ctx, client.ReadRequest{
		Filename: responseFilename,
	})

	if err != nil {
		return nil, err
	}

	respHeader := http.Header{}

	// If there was a single Accept request header, it is the manifest media
	// type. Add it to the response as Content-Type.
	if accept := parseSingleAccept(req.Header.Get("Accept")); accept != "" {
		respHeader.Set("Content-Type", accept)
	}

	return &http.Response{
		Status:        "200 OK",
		StatusCode:    200,
		Proto:         req.Proto,
		ProtoMajor:    req.ProtoMajor,
		ProtoMinor:    req.ProtoMinor,
		Header:        respHeader,
		ContentLength: int64(len(dt)),
		Body:          io.NopCloser(bytes.NewReader(dt)),
		Uncompressed:  true,
		Request:       req,
	}, nil
}

func (rt *registryTransport) requestToLLB(req *http.Request) llb.State {
	opts := append(
		rt.options,
		llb.Header(llb.HTTPHeader{
			Accept:    req.Header.Get("Accept"),
			UserAgent: req.Header.Get("User-Agent"),
		}),
		llb.Filename(responseFilename),
		llb.WithCustomNamef("🎱 fetch %s", req.URL),
	)

	if secretName, ok := rt.buildArgs[authSecretNameOptionPrefix+req.URL.Host]; ok {
		opts = append(opts, llb.AuthHeaderSecret(secretName))
	}

	return llb.HTTP(req.URL.String(), opts...)
}

func parseSingleAccept(accept string) string {
	parts := strings.SplitN(accept, ",", 1)
	if len(parts) != 1 {
		return ""
	}

	parts = strings.SplitN(parts[0], ";", 1)
	if len(parts) < 1 {
		return ""
	}

	singleAccept := strings.Trim(parts[0], " ")

	if strings.HasPrefix(singleAccept, "*") || strings.HasSuffix(singleAccept, "*") {
		return ""
	}

	return singleAccept
}
