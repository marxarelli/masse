package v1

import (
	"os"

	"github.com/moby/buildkit/client/llb"
)

func (c *compiler) Expand(state llb.State, s string) string {
	return os.Expand(s, func(name string) string {
		replacement, found, err := state.GetEnv(c.ctx, name)
		if err == nil && found {
			return replacement
		}

		return ""
	})
}
