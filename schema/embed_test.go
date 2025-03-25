package schema

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFS(t *testing.T) {
	req := require.New(t)

	expectedFiles := []string{
		"config.cue",
		"common/creation.cue",
		"state/run.cue",
		"target/targets.cue",
	}

	for _, expectedFile := range expectedFiles {
		_, err := fs.Stat(FS, expectedFile)
		req.NoError(err)
	}
}
