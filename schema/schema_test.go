package schema

import (
	"testing"

	"cuelang.org/go/cue/cuecontext"
	"github.com/stretchr/testify/require"
)

func TestLoaderConfig(t *testing.T) {
	req := require.New(t)
	ctx := cuecontext.New()

	cfg, err := LoaderConfig(ctx, "/root/dir")

	req.NoError(err)
	req.Equal("/", cfg.Dir)

	overlayEntries := make([]string, len(cfg.Overlay))
	i := 0
	for entry := range cfg.Overlay {
		overlayEntries[i] = entry
		i++
	}

	req.Contains(overlayEntries, "/root/dir/cue.mod/pkg/wikimedia.org/dduvall/phyton/common/creation.cue")
	req.Contains(overlayEntries, "/root/dir/cue.mod/pkg/wikimedia.org/dduvall/phyton/source/source.cue")
	req.Contains(overlayEntries, "/root/dir/cue.mod/pkg/wikimedia.org/dduvall/phyton/op/op.cue")
	req.Contains(overlayEntries, "/root/dir/cue.mod/pkg/wikimedia.org/dduvall/phyton/state/state.cue")
}
