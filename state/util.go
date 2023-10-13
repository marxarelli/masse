package state

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/moby/buildkit/client/llb"
)

func qualifyStatePath(state llb.State, path string) string {
	if filepath.IsAbs(path) {
		return path
	}

	hadTrailingSlash := strings.HasSuffix(path, "/")

	cwd, _ := state.GetDir(context.TODO())
	abs := filepath.Join(cwd, path)

	if hadTrailingSlash {
		abs += "/"
	}

	return abs
}
