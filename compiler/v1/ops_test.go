package v1

import (
	"testing"

	"gitlab.wikimedia.org/dduvall/masse/util/llbtest"
	"gitlab.wikimedia.org/dduvall/masse/util/testcompile"
)

func TestOps(t *testing.T) {
	compile := testcompile.New(
		t,
		[]string{"wikimedia.org/dduvall/masse/state"},
		testcompile.WithCompiler(func() *compiler {
			return newCompiler(nil)
		}),
	)

	compile.Test(
		"ops",
		`state.#Ops & { ops: [{ run: "make libs" }, { run: "make bin" }] }`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			_, eops := req.ContainsNExecOps(2)
			req.Equal([]string{"/bin/sh", "-c", "make libs"}, eops[0].Exec.Meta.Args)
			req.Equal([]string{"/bin/sh", "-c", "make bin"}, eops[1].Exec.Meta.Args)
		},
	)
}