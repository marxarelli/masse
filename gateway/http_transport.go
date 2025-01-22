package gateway

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/gateway/client"
	"github.com/pkg/errors"
)

const (
	authSecretNameOptionPrefix = "CUE_REGISTRY_AUTH_SECRET."
	responseFilename           = "/http.response"
)

var (
	resourceIcons = map[string]string{
		"manifests": "📦📋",
		"blobs":     "📦💾",
	}
)

type httpTransport struct {
	ctx       context.Context
	client    client.Client
	buildArgs map[string]string
}

func NewRegistryTransport(gw *Gateway) http.RoundTripper {
	return &httpTransport{
		ctx:       gw.ctx,
		client:    gw.bkClient,
		buildArgs: gw.Config.BuildArgs,
	}
}

func (ht *httpTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Method != "GET" {
		return nil, errors.New("can only perform GET requests")
	}

	def, err := ht.requestToLLB(req).Marshal(ht.ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal HTTP LLB state")
	}

	res, err := ht.client.Solve(ht.ctx, client.SolveRequest{
		Definition: def.ToPB(),
	})
	if err != nil {
		return nil, err
	}

	ref, err := res.SingleRef()
	if err != nil {
		return nil, err
	}

	dt, err := ref.ReadFile(ht.ctx, client.ReadRequest{
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

func (ht *httpTransport) requestToLLB(req *http.Request) llb.State {
	opts := []llb.HTTPOption{
		llb.Header(llb.HTTPHeader{
			Accept:    req.Header.Get("Accept"),
			UserAgent: req.Header.Get("User-Agent"),
		}),
		llb.Filename(responseFilename),
	}

	// Create a custom name based on the registry request path, either of:
	// /v2/{repo}/manifests/{tag}
	// /v2/{repo}/blobs/{sha}
	if rPath, found := strings.CutPrefix(req.URL.Path, "/v2/"); found {
		p := strings.Split(rPath, "/")
		l := len(p)
		if l > 2 {
			repo := path.Join(p[0 : l-2]...)
			kind := p[l-2]
			tag := p[l-1]

			opts = append(
				opts,
				llb.WithCustomNamef(
					"%s load CUE module %s/%s:%s",
					resourceIcons[kind], req.URL.Host, repo, tag,
				),
			)
		}
	}

	if secretName, ok := ht.buildArgs[authSecretNameOptionPrefix+req.URL.Host]; ok {
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
