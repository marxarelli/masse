package v1

import (
	"testing"

	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/masse/util/llbtest"
	"gitlab.wikimedia.org/dduvall/masse/util/testcompile"
)

func TestDiff(t *testing.T) {
	compile := testcompile.New(
		t,
		[]string{"wikimedia.org/dduvall/masse/state"},
		testcompile.WithCompiler(func() *compiler {
			return newCompiler(nil)
		}),
	)

	compile.Test(
		"chain",
		`state.#Op & { diff: { run: ["make", "libs"] } }`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			ops, dops := req.ContainsNDiffOps(1)
			req.Equal(int64(0), dops[0].Diff.Lower.Input)
			req.Equal(int64(1), dops[0].Diff.Upper.Input)

			inputs := req.HasValidInputs(ops[0])
			req.Len(inputs, 2)

			_, sops := req.ContainsNSourceOps(1)
			_, eops := req.ContainsNExecOps(1)
			req.Equal(inputs[0].Op, sops[0])
			req.Equal(inputs[1].Op, eops[0])
		},
		testcompile.WithInitialState(llb.Git("https://an.example/repo.git", "refs/heads/main")),
	)

	compile.Test(
		"from",
		`state.#Op & { diff: { run: ["make", "libs"] }, from: "repo" }`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			ops, dops := req.ContainsNDiffOps(1)
			req.Equal(int64(0), dops[0].Diff.Lower.Input)
			req.Equal(int64(1), dops[0].Diff.Upper.Input)

			inputs := req.HasValidInputs(ops[0])
			req.Len(inputs, 2)

			_, sops := req.ContainsNSourceOps(1)
			_, eops := req.ContainsNExecOps(1)
			req.Equal(inputs[0].Op, sops[0])
			req.Equal(inputs[1].Op, eops[0])
		},
		testcompile.WithCompiler(func() *compiler {
			c := newCompiler(nil)
			c.chainCompilers = map[string]chainCompiler{
				"repo": func(_ *compiler) *chainResult {
					return &chainResult{state: llb.Git("https://an.example/repo.git", "refs/heads/main")}
				},
			}
			return c
		}),
	)

	compile.Test(
		"options",
		`state.#Op & { diff: { run: ["make", "libs"] }, options: customName: "libs" }`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			ops, dops := req.ContainsNDiffOps(1)
			req.Equal(int64(0), dops[0].Diff.Lower.Input)
			req.Equal(int64(1), dops[0].Diff.Upper.Input)

			inputs := req.HasValidInputs(ops[0])
			req.Len(inputs, 2)

			_, sops := req.ContainsNSourceOps(1)
			_, eops := req.ContainsNExecOps(1)
			req.Equal(inputs[0].Op, sops[0])
			req.Equal(inputs[1].Op, eops[0])
		},
		testcompile.WithInitialState(llb.Git("https://an.example/repo.git", "refs/heads/main")),
	)
}
