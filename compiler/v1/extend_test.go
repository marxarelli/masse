package v1

import (
	"testing"

	"github.com/moby/buildkit/client/llb"
	oci "github.com/opencontainers/image-spec/specs-go/v1"
	"gitlab.wikimedia.org/dduvall/masse/util/testcompile"
	"gitlab.wikimedia.org/dduvall/masse/util/testmetaresolver"
	"gitlab.wikimedia.org/repos/releng/llbtest/llbtest"
)

func TestExtend(t *testing.T) {
	compile := testcompile.New(
		t,
		[]string{"github.com/marxarelli/masse/state"},
		testcompile.WithCompiler(func() *compiler {
			c := newCompiler(nil, WithImageMetaResolver(testmetaresolver.New("golang:1.23", oci.Image{})))
			c.chainCompilers = map[string]chainCompiler{
				"go": func(_ *compiler) *chainResult {
					return &chainResult{
						state: llb.Image("golang:1.23"),
					}
				},
			}
			return c
		}),
	)

	compile.Test(
		"minimal",
		`state.#Extend & { extend: "go" }`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			_, sops := req.ContainsNSourceOps(1)
			req.Equal("docker-image://docker.io/library/golang:1.23", sops[0].Source.Identifier)
		},
	)
}
