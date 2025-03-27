package v1

import (
	"context"
	"testing"

	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/masse/util/testcompile"
	"gitlab.wikimedia.org/repos/releng/llbtest/llbtest"
)

func TestContext(t *testing.T) {
	compile := testcompile.New(
		t,
		[]string{"github.com/marxarelli/masse/state"},
	)

	compile.Test(
		"main context",
		`state.#Source & { mainContext: true }`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			_, sops := req.ContainsNSourceOps(1)
			req.Equal("local://foo", sops[0].Source.Identifier)
		},
		testcompile.WithCompiler(func() *compiler {
			compiler := newCompiler(
				nil,
				WithMainContextLoader(func(_ context.Context, opts ...llb.LocalOption) (*llb.State, error) {
					st := llb.Local("foo", opts...)
					return &st, nil
				}),
			)
			return compiler
		}),
	)

	compile.Test(
		"named context",
		`state.#Source & { context: "foo" }`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			_, sops := req.ContainsNSourceOps(1)
			req.Equal("local://foo", sops[0].Source.Identifier)
		},
		testcompile.WithCompiler(func() *compiler {
			compiler := newCompiler(
				nil,
				WithNamedContextLoader(func(_ context.Context, name string, opts ...llb.LocalOption) (*llb.State, error) {
					st := llb.Local(name, opts...)
					return &st, nil
				}),
			)
			return compiler
		}),
	)
}
