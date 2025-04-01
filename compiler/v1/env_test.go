package v1

import (
	"testing"

	"gitlab.wikimedia.org/dduvall/masse/util/testcompile"
	"gitlab.wikimedia.org/repos/releng/llbtest/llbtest"
)

func TestEnv(t *testing.T) {
	compile := testcompile.New(
		t,
		[]string{"github.com/marxarelli/masse/state"},
		testcompile.WithCompiler(func() *compiler {
			return newCompiler(nil)
		}),
	)

	compile.Test(
		"env",
		`state.#Op & { run: ["make"], options: env: FOO: "foo$bar" }`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			_, eops := req.ContainsNExecOps(1)
			req.Equal([]string{"FOO=foo$bar"}, eops[0].Exec.Meta.Env)
		},
	)

	compile.Test(
		"expandVariables",
		`state.#Op & { run: ["make"], options: [ { env: BAR: "bar" }, { env: FOO: "foo${BAR}", expandVariables: true } ] }`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			_, eops := req.ContainsNExecOps(1)
			req.Equal([]string{"BAR=bar", "FOO=foobar"}, eops[0].Exec.Meta.Env)
		},
	)
}
