package schema

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFS(t *testing.T) {
	req := require.New(t)

	_, err := fs.ReadDir(FS, "target")
	req.NoError(err)
}
