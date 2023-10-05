package load

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"cuelang.org/go/cue"
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

func lines(lns ...string) string {
	return strings.Join(lns, "\n") + "\n"
}
