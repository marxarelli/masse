package testconfig

import (
	"strings"
	"testing"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"github.com/stretchr/testify/require"
	"gitlab.wikimedia.org/dduvall/masse/load"
)

type Tester struct {
	*testing.T
	CUEImports []string
}

func New(t *testing.T, imports []string) *Tester {
	return &Tester{T: t, CUEImports: imports}
}

func (tester *Tester) Run(name string, f func(*Tester)) {
	tester.Helper()

	tester.T.Run(name, func(t *testing.T) {
		f(New(t, tester.CUEImports))
	})
}

func (tester *Tester) Test(expr string, f func(*testing.T, *require.Assertions, cue.Value)) {
	tester.Helper()

	tester.T.Run(expr, func(t *testing.T) {
		t.Helper()
		t.Parallel()
		req := require.New(t)

		ctx := cuecontext.New(cuecontext.EvaluatorVersion(cuecontext.EvalV2))

		dir := t.TempDir()

		main, err := load.MainInstanceWith(
			dir,
			map[string][]byte{
				"test.cue": []byte(
					"package main\n" +
						"import (\n  \"" +
						strings.Join(tester.CUEImports, "\"\n  \"") +
						"\"\n)\n" +
						`x: ` + expr,
				),
			},
		)

		require.NoError(t, err)

		inst := ctx.BuildInstance(main)
		require.NoError(t, inst.Err())

		x := inst.LookupPath(cue.ParsePath("x"))
		require.NoError(t, x.Err())

		f(t, req, x)
	})
}
