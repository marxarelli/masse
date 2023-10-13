package load

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	req := require.New(t)

	tmpdir, err := os.MkdirTemp("", "phyton-load-test")
	defer os.RemoveAll(tmpdir)
	req.NoError(err)

	cueFile := filepath.Join(tmpdir, "test.cue")
	err = os.WriteFile(cueFile, []byte(lines(
		`package main`,
		``,
		`import "wikimedia.org/dduvall/phyton/schema/state"`,
		``,
		`foo: state.#Run & { run: "foo" }`,
	)), 0644)
	req.NoError(err)

	val, err := LoadPath(cueFile)
	req.NoError(err)
	req.NoError(val.Err())

	foo := val.LookupPath(cue.ParsePath("foo"))
	req.NoError(foo.Err())
}

func TestMainInstanceWith(t *testing.T) {
	req := require.New(t)

	main, err := MainInstanceWith(
		map[string][]byte{
			"/foo.cue": []byte(lines(
				`chains: {`,
				`  foo: [`,
				`    { image: "foo.example/image/ref" },`,
				`  ]`,
				`}`,
				`layouts: foo: comprises: ["foo"]`,
			)),
		},
	)
	req.NotNil(main)
	req.NoError(err)

	ctx := cuecontext.New()
	value := ctx.BuildInstance(main)
	req.NoError(value.Err())

	chains := value.LookupPath(cue.ParsePath("chains"))
	req.NoError(chains.Err())

	ref := value.LookupPath(cue.ParsePath("chains.foo[0].image"))
	req.NoError(ref.Err())

	req.Equal(cue.StringKind, ref.Kind())
	refValue, err := ref.String()
	req.NoError(err)
	req.Equal("foo.example/image/ref", refValue)
}

func lines(lns ...string) string {
	return strings.Join(lns, "\n") + "\n"
}
