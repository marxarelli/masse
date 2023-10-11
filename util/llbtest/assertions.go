package llbtest

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/solver/pb"
	digest "github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/require"
)

// Assertions provides methods for making assertions on LLB build graphs,
// specifically the marshaled [llb.Definition] and [pb.Op] therein. It extends
// [require.Assertions] to provide a singular assertions interface in LLB
// related tests.
type Assertions struct {
	*require.Assertions

	t     *testing.T
	def   *llb.Definition
	opMap map[digest.Digest]pb.Op
	ops   []pb.Op
}

// New returns a new [Assertions] that encapsulates the [testing.T] and
// [llb.Definition] for ease of making many assertions about the contained op.
func New(t *testing.T, def *llb.Definition) *Assertions {
	t.Helper()

	opMap, ops := ParseDef(t, def.Def)

	return &Assertions{
		Assertions: require.New(t),
		t:          t,
		def:        def,
		opMap:      opMap,
		ops:        ops,
	}
}

// ContainsNOps requires that there are exactly n [pb.Op] that are of the
// given type. If successful, the matching [pb.Op] and encapsulated [pb.Op.Op]
// field of the given type are returned.
func ContainsNOps[T any](t *testing.T, ops []pb.Op, n int, msg string) ([]pb.Op, []T) {
	t.Helper()

	matchingOps := []pb.Op{}
	matchingTOps := []T{}

	for _, op := range ops {
		tOp, ok := op.Op.(T)

		if ok {
			matchingOps = append(matchingOps, op)
			matchingTOps = append(matchingTOps, tOp)
		}
	}

	require.Lenf(t, matchingTOps, n, msg, n, spew.Sdump(ops))

	return matchingOps, matchingTOps
}

// ContainsNFileActions requires that there are exactly n [pb.FileAction]
// within the [pb.Op_File] that are of the given type. If successful, the
// matching [pb.FileAction] and encapsulated [pb.FileAction.Action] field of
// the given type are returned.
func ContainsNFileActions[T any](t *testing.T, fileOp *pb.Op_File, n int, msg string) ([]*pb.FileAction, []T) {
	t.Helper()

	matchingActions := []*pb.FileAction{}
	matchingTActions := []T{}

	if fileOp != nil {
		for _, action := range fileOp.File.Actions {
			tAction, ok := action.Action.(T)

			if ok {
				matchingActions = append(matchingActions, action)
				matchingTActions = append(matchingTActions, tAction)
			}
		}
	}

	require.Lenf(t, matchingTActions, n, msg, n, spew.Sdump(fileOp.File.Actions))

	return matchingActions, matchingTActions
}

// HasValidInputs takes a [pb.Op] and asserts that all of its [pb.Op.Inputs]
// have a corresponding [pb.Op] in the given op map. If the assertion
// succeeds, each corresponding [pb.Op] is returned.
func HasValidInputs(t *testing.T, opMap map[digest.Digest]pb.Op, op pb.Op) []pb.Op {
	t.Helper()

	inputOps := make([]pb.Op, len(op.Inputs))

	for i, input := range op.Inputs {
		inputOp, ok := opMap[input.Digest]

		require.True(
			t,
			ok,
			"op %v has an input (%+v) that was not found in the definition %+v",
			op, input, opMap,
		)

		inputOps[i] = inputOp
	}

	return inputOps
}

// ContainsNSourceOps is a convenience method for ContainsNOps[pb.Op_Source](...)
func (llbt *Assertions) ContainsNSourceOps(n int) ([]pb.Op, []*pb.Op_Source) {
	return ContainsNOps[*pb.Op_Source](llbt.t, llbt.ops, n, "should contain %d source ops: %s")
}

// ContainsNFileOps is a convenience method for ContainsNOps[pb.Op_File](...)
func (llbt *Assertions) ContainsNFileOps(n int) ([]pb.Op, []*pb.Op_File) {
	return ContainsNOps[*pb.Op_File](llbt.t, llbt.ops, n, "should contain %d file ops: %s")
}

// ContainsNExecOps is a convenience method for ContainsNOps[pb.Op_Exec](...)
func (llbt *Assertions) ContainsNExecOps(n int) ([]pb.Op, []*pb.Op_Exec) {
	return ContainsNOps[*pb.Op_Exec](llbt.t, llbt.ops, n, "should contain %d exec ops: %s")
}

// ContainsNCopyActions is a convenience method for
// ContainsNFileActions[pb.FileAction_Copy](...)
func (llbt *Assertions) ContainsNCopyActions(fileOp *pb.Op_File, n int) ([]*pb.FileAction, []*pb.FileAction_Copy) {
	return ContainsNFileActions[*pb.FileAction_Copy](llbt.t, fileOp, n, "should contain %d copy actions: %s")
}

// HasValidInputs is a convenience method for [HasValidInputs]
func (llbt *Assertions) HasValidInputs(op pb.Op) []pb.Op {
	return HasValidInputs(llbt.t, llbt.opMap, op)
}
