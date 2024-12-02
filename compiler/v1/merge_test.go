package v1

import (
	"testing"

	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/masse/util/llbtest"
	"gitlab.wikimedia.org/dduvall/masse/util/testcompile"
)

func TestMerge(t *testing.T) {
	compile := testcompile.New(
		t,
		[]string{"wikimedia.org/dduvall/masse/state"},
		testcompile.WithCompiler(func() *compiler {
			c := newCompiler(nil)
			c.chainCompilers = map[string]chainCompiler{
				"repo": func(_ *compiler) *chainResult {
					return &chainResult{state: llb.Git("an.example/repo.git", "refs/heads/main")}
				},
				"files": func(_ *compiler) *chainResult {
					return &chainResult{state: llb.Scratch().File(llb.Copy(llb.Local("context"), "/src", "/dest"))}
				},
			}
			return c
		}),
	)

	compile.Test(
		"minimal",
		`state.#Merge & { merge: ["repo", "files"] }`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			ops, mops := req.ContainsNMergeOps(1)
			req.Len(mops[0].Merge.Inputs, 2)

			_, sops := req.ContainsNSourceOps(2)
			_, fops := req.ContainsNFileOps(1)

			iops := req.HasValidInputs(ops[0])
			req.Len(iops, 2)
			req.Equal(iops[0].Op, sops[0])
			req.Equal(iops[1].Op, fops[0])
		},
	)
}
