package v1

import (
	"testing"

	"github.com/moby/buildkit/client/llb"
	oci "github.com/opencontainers/image-spec/specs-go/v1"
	"gitlab.wikimedia.org/dduvall/masse/target"
	"gitlab.wikimedia.org/dduvall/masse/util/llbtest"
	"gitlab.wikimedia.org/dduvall/masse/util/testcompile"
	"gitlab.wikimedia.org/dduvall/masse/util/testmetaresolver"
)

func TestTarget(t *testing.T) {
	compile := testcompile.New(
		t,
		[]string{"wikimedia.org/releng/masse/target"},
		testcompile.WithCompiler(func() *compiler {
			c := newCompiler(nil)
			c.chainCompilers = map[string]chainCompiler{
				"repo": func(_ *compiler) *chainResult {
					return &chainResult{
						state: llb.Git("an.example/repo.git", "refs/heads/main"),
					}
				},
				"go": func(_ *compiler) *chainResult {
					return &chainResult{
						state: llb.Image("golang:1.23", llb.WithMetaResolver(testmetaresolver.New("golang:1.23", oci.Image{}))),
					}
				},
			}
			return c
		}),
		testcompile.WithCompileFunc(func(test *testcompile.Test) (llb.State, error) {
			target := &target.Target{}

			err := target.UnmarshalCUE(test.Value)
			if err != nil {
				return test.State, err
			}

			return test.Compiler.Compile(target)
		}),
	)

	compile.Test(
		"target",
		`target.#Target & {
			build: [
				{ extend: "go" },
				{ file: [
					{ copy: "go.mod", from: "repo" },
					{ copy: "go.sum", from: "repo" },
				] },
				{ sh: "go mod download" },
				{ file: { copy: ".", from: "repo" } },
				{ sh: "go build ./cmd/foo"
					options: [ { cache: "/root/.cache/go-build", access: "locked" } ] },
			],
		}`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			_, sops := req.ContainsNSourceOps(2)
			req.Equal("docker-image://docker.io/library/golang:1.23", sops[0].Source.Identifier)
			req.Equal("git://an.example/repo.git#refs/heads/main", sops[1].Source.Identifier)
		},
	)
}
