package load

import (
	"strings"
	"testing"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"github.com/stretchr/testify/require"
)

func TestMainInstanceWith(t *testing.T) {
	req := require.New(t)
	dir := t.TempDir()

	main, err := MainInstanceWith(
		dir,
		map[string][]byte{
			"foo.cue": []byte(lines(
				`package main`,
				``,
				`chains: {`,
				`  foo: [`,
				`    { image: "foo.example/image/ref" },`,
				`  ]`,
				`}`,
				``,
				`targets: foo: build: "foo"`,
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
