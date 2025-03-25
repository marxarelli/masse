package schema

import (
	"path"
	"testing"

	"cuelang.org/go/cue/load"
	"github.com/stretchr/testify/require"
)

func TestNewModuleFS(t *testing.T) {
	req := require.New(t)
	dir := t.TempDir()
	expectedEmbedDir := path.Join(dir, "cue.mod/embed/masse")

	cfg := &load.Config{Dir: dir}
	mfs, err := NewModuleFS(cfg)
	req.NoError(err)

	req.Equal(expectedEmbedDir, mfs.OSRoot())
	req.Contains(cfg.Overlay, path.Join(expectedEmbedDir, "config.cue"))
}
