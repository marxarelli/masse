package schema

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoaderConfig(t *testing.T) {
	req := require.New(t)
	dir := t.TempDir()

	cfg, err := LoaderConfig(dir)

	req.NoError(err)
	req.Equal(dir, cfg.Dir)

	overlayEntries := make([]string, len(cfg.Overlay))
	i := 0
	for entry := range cfg.Overlay {
		overlayEntries[i] = entry
		i++
	}

	req.Contains(overlayEntries, filepath.Join(dir, "/cue.mod/pkg/github.com/marxarelli/masse/common/creation.cue"))
	req.Contains(overlayEntries, filepath.Join(dir, "/cue.mod/pkg/github.com/marxarelli/masse/state/state.cue"))
}
