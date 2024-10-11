package llbtest

import (
	"context"
	"testing"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/solver/pb"
	digest "github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/require"
)

// ParseState marshals the given [llb.State] to a [llb.Definition], then
// parses it as a [pb.Op] slice and map of digests to ops.
func ParseState(t *testing.T, state llb.State) (map[string]pb.Op, []pb.Op) {
	t.Helper()

	def, err := state.Marshal(context.TODO())
	require.NoError(t, err)

	return ParseDef(t, def.Def)
}

// ParseDef parses the given [llb.Definition] in a [pb.Op] slice and map of
// digests to ops.
func ParseDef(t *testing.T, def [][]byte) (map[string]pb.Op, []pb.Op) {
	t.Helper()

	m := map[string]pb.Op{}
	arr := make([]pb.Op, 0, len(def))

	for _, dt := range def {
		var op pb.Op
		err := (&op).Unmarshal(dt)
		require.NoError(t, err)
		dgst := digest.FromBytes(dt)
		m[string(dgst)] = op
		arr = append(arr, op)
		// fmt.Printf(":: %T %+v\n", op.Op, op)
	}

	return m, arr
}

// LastOp returns the final Op in the given ordered slice
func LastOp(t *testing.T, arr []pb.Op) (string, int) {
	t.Helper()

	require.True(t, len(arr) > 1)

	op := arr[len(arr)-1]
	require.Equal(t, 1, len(op.Inputs))
	return op.Inputs[0].Digest, int(op.Inputs[0].Index)
}
