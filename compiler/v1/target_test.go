package v1

import (
	"testing"

	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/masse/target"
	"gitlab.wikimedia.org/dduvall/masse/util/llbtest"
	"gitlab.wikimedia.org/dduvall/masse/util/testcompile"
)

func TestTarget(t *testing.T) {
	compile := testcompile.New(
		t,
		[]string{"github.com/marxarelli/masse/target"},
		testcompile.WithCompiler(func() *compiler {
			c := newCompiler(nil)
			c.chainCompilers = map[string]chainCompiler{
				"foo": func(_ *compiler) *chainResult {
					return &chainResult{
						state: llb.Git("an.example/repo.git", "refs/heads/main"),
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
			build: "foo"
		}`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			req.ContainsNSourceOps(1)
		},
	)
}
