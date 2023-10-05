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
	req.Equal("/root/dir", cfg.Dir)

	overlayEntries := make([]string, len(cfg.Overlay))
	i := 0
	for entry := range cfg.Overlay {
		overlayEntries[i] = entry
		i++
	}

	req.Contains(overlayEntries, "/root/dir/cue.mod/pkg/wikimedia.org/dduvall/phyton/schema/common/creation.cue")
	req.Contains(overlayEntries, "/root/dir/cue.mod/pkg/wikimedia.org/dduvall/phyton/schema/state/state.cue")
}

func TestInstances(t *testing.T) {
	req := require.New(t)
	ctx := cuecontext.New()

	instances, err := Instances(ctx)
	req.NoError(err)
	req.Len(instances, 4)

	for _, ins := range instances {
		val := ins.Value()
		req.NotNil(val)
		req.NoError(val.Err())
	}
}
