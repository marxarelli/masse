package v1

import (
	"testing"

	"cuelang.org/go/cue"
	"github.com/stretchr/testify/require"
	"gitlab.wikimedia.org/dduvall/masse/util/testconfig"
)

func TestCompilerCycleDetection(t *testing.T) {
	tester := testconfig.New(
		t,
		[]string{"wikimedia.org/dduvall/masse/state"},
	)

	tester.Test(
		"chains/cycles",
		`state.#Chains & { foo: [ { extend: "bar" } ], bar: [ { extend: "foo" } ] }`,
		func(_ *testing.T, req *require.Assertions, v cue.Value) {
			chains := map[string]cue.Value{}
			req.NoError(v.Decode(&chains))
			compiler := newCompiler(chains)
			compiler.defineChainCompilers()
			_, err := compiler.compileChainByRef(v.Context().CompileString(`"foo"`))
			req.ErrorContains(err, "chain ref cycle detected")
		},
	)
}
