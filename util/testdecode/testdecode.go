package testdecode

import (
	"fmt"
	"strings"
	"testing"

	"cuelang.org/go/cue"
	"github.com/stretchr/testify/require"
	"gitlab.wikimedia.org/dduvall/masse/load"
	"gitlab.wikimedia.org/dduvall/masse/schema"
)

type Tester struct {
	*testing.T
	CUEImports []string
}

func (t *Tester) LoadCUE(code ...string) cue.Value {
	t.Helper()

	cueImports := make([]string, len(t.CUEImports))
	for i, imp := range t.CUEImports {
		cueImports[i] = fmt.Sprintf(`  %q`, imp)
	}

	cuePackage := append(
		[]string{
			`package main`,
			``,
			`import (`,
		},
		cueImports...,
	)
	cuePackage = append(cuePackage, ")")
	cuePackage = append(cuePackage, code...)

	val, err := load.LoadBytes([]byte(strings.Join(cuePackage, "\n")+"\n"), "/config.cue")
	require.NoError(t, err)
	require.NoError(t, val.Err())

	return val
}

func Run[T any](tester *Tester, name, cueCode string, expected T) {
	tester.Helper()

	tester.Run(name, func(t *testing.T) {
		t.Helper()
		t.Parallel()

		req := require.New(t)

		actual := new(T)
		err := schema.Decode(tester.LoadCUE(cueCode), actual)
		req.NoError(err)

		req.Equal(&expected, actual)
	})
}
