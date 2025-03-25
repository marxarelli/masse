package v1

import (
	"testing"

	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/masse/util/testcompile"
	"gitlab.wikimedia.org/repos/releng/llbtest/llbtest"
)

func TestIgnoreCache(t *testing.T) {
	compile := testcompile.New(
		t,
		[]string{"github.com/marxarelli/masse/state"},
		testcompile.WithCompiler(func() *compiler {
			c := newCompiler(nil, WithIgnoreCache(true))
			c.chainCompilers = map[string]chainCompiler{
				"local": func(_ *compiler) *chainResult {
					return &chainResult{state: llb.Local("context").Dir("/src")}
				},
			}
			return c
		}),
	)

	compile.Run("file", func(compile *testcompile.Tester) {
		compile.Test(
			"file",
			`state.#File & { file: { copy: "./foo", from: "local" } }`,
			func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
				ops, _ := req.ContainsNFileOps(1)
				md := req.HasMetadata(ops[0])
				req.True(md.IgnoreCache)
			},
		)
	})
}
