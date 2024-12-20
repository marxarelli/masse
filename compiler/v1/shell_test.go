package v1

import (
	"testing"

	"gitlab.wikimedia.org/dduvall/masse/util/llbtest"
	"gitlab.wikimedia.org/dduvall/masse/util/testcompile"
)

func TestShell(t *testing.T) {
	compile := testcompile.New(
		t,
		[]string{"wikimedia.org/releng/masse/state"},
		testcompile.WithCompiler(func() *compiler {
			return newCompiler(nil)
		}),
	)

	compile.Test(
		"minimal",
		`state.#Op & { sh: "make" }`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			_, eops := req.ContainsNExecOps(1)
			req.Equal([]string{"/bin/sh", "-c", "make"}, eops[0].Exec.Meta.Args)
		},
	)

	compile.Test(
		"arguments/single",
		`state.#Op & { sh: "make", arguments: "foo bar" }`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			_, eops := req.ContainsNExecOps(1)
			req.Equal([]string{"/bin/sh", "-c", `make "foo bar"`}, eops[0].Exec.Meta.Args)
		},
	)

	compile.Test(
		"arguments/multiple",
		`state.#Op & { sh: "make", arguments: ["foo", "bar"] }`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			_, eops := req.ContainsNExecOps(1)
			req.Equal([]string{"/bin/sh", "-c", `make "foo" "bar"`}, eops[0].Exec.Meta.Args)
		},
	)

	compile.Test(
		"customName",
		`state.#Op & { sh: "make" }`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			ops, _ := req.ContainsNExecOps(1)
			md := req.HasMetadata(ops[0])
			req.Contains(md.Description, "llb.customname")
			req.Equal("💻 make", md.Description["llb.customname"])
		},
	)

	compile.Test(
		"customName/arguments",
		`state.#Op & { sh: "make", arguments: ["foo"] }`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			ops, _ := req.ContainsNExecOps(1)
			md := req.HasMetadata(ops[0])
			req.Contains(md.Description, "llb.customname")
			req.Equal(`💻 make "foo"`, md.Description["llb.customname"])
		},
	)
}
