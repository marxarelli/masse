package v1

import (
	"testing"

	"gitlab.wikimedia.org/dduvall/masse/util/llbtest"
	"gitlab.wikimedia.org/dduvall/masse/util/testcompile"
)

func TestLocal(t *testing.T) {
	compile := testcompile.New(
		t,
		[]string{"wikimedia.org/dduvall/masse/state"},
		testcompile.WithCompiler(func() *compiler {
			return newCompiler(nil)
		}),
	)

	compile.Test(
		"minimal",
		`state.#Local & { local: "foo" }`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			_, sops := req.ContainsNSourceOps(1)
			req.Equal("local://foo", sops[0].Source.Identifier)
		},
	)

	compile.Test(
		"options/include",
		`state.#Local & { local: "foo", options: include: ["*.c"] }`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			_, sops := req.ContainsNSourceOps(1)
			req.Equal("local://foo", sops[0].Source.Identifier)
			req.Contains(sops[0].Source.Attrs, "local.includepattern")
			req.Equal(`["*.c"]`, sops[0].Source.Attrs["local.includepattern"])
		},
	)
}
