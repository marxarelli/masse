package schema

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoaderConfig(t *testing.T) {
	req := require.New(t)

	cfg, err := LoaderConfig("/root/dir")

	req.NoError(err)
	req.Equal("/root/dir", cfg.Dir)

	overlayEntries := make([]string, len(cfg.Overlay))
	i := 0
	for entry := range cfg.Overlay {
		overlayEntries[i] = entry
		i++
	}

	req.Contains(overlayEntries, "/root/dir/cue.mod/pkg/wikimedia.org/dduvall/masse/common/creation.cue")
	req.Contains(overlayEntries, "/root/dir/cue.mod/pkg/wikimedia.org/dduvall/masse/state/state.cue")
}
