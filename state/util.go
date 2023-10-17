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

func constraintsTo[T any](opts []llb.ConstraintsOpt) []T {
	ts := make([]T, len(opts))
	for i, x := range opts {
		ts[i] = x.(T)
	}
	return ts
}
