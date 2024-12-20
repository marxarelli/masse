package v1

import (
	"testing"

	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/masse/util/llbtest"
	"gitlab.wikimedia.org/dduvall/masse/util/testcompile"
)

func TestCopyOp(t *testing.T) {
	compile := testcompile.New(
		t,
		[]string{"github.com/marxarelli/masse/state"},
		testcompile.WithCompiler(func() *compiler {
			c := newCompiler(nil)
			c.chainCompilers = map[string]chainCompiler{
				"local": func(_ *compiler) *chainResult {
					return &chainResult{state: llb.Local("context").Dir("/src")}
				},
			}
			return c
		}),
	)

	compile.Run("CopyOp", func(compile *testcompile.Tester) {
		compile.Test(
			"single",
			`state.#Op & { copy: "./foo", from: "local" }`,
			func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
				_, fops := req.ContainsNFileOps(1)
				_, copies := req.ContainsNCopyActions(fops[0], 1)

				req.Equal("/src/foo", copies[0].Copy.Src)
				req.Equal("/", copies[0].Copy.Dest)
				req.False(copies[0].Copy.DirCopyContents)
				req.False(copies[0].Copy.FollowSymlink)
			},
		)

		compile.Test(
			"multiple",
			`state.#Op & { copy: ["./foo", "./bar"], from: "local" }`,
			func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
				_, fops := req.ContainsNFileOps(1)
				_, copies := req.ContainsNCopyActions(fops[0], 2)

				req.Equal("/src/foo", copies[0].Copy.Src)
				req.Equal("/", copies[0].Copy.Dest)
				req.False(copies[0].Copy.DirCopyContents)
				req.False(copies[0].Copy.FollowSymlink)

				req.Equal("/src/bar", copies[1].Copy.Src)
				req.Equal("/", copies[1].Copy.Dest)
				req.False(copies[1].Copy.DirCopyContents)
				req.False(copies[1].Copy.FollowSymlink)
			},
		)
	})
}
