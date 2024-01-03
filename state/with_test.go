package state

import (
	"context"
	"testing"

	"github.com/moby/buildkit/client/llb"
	"github.com/stretchr/testify/require"
	"gitlab.wikimedia.org/dduvall/masse/common"
	"gitlab.wikimedia.org/dduvall/masse/util/llbtest"
)

func TestCompileWith(t *testing.T) {
	req := require.New(t)

	t.Run("Env", func(t *testing.T) {
		with := &With{
			With: Options{
				{Env: &Env{Env: common.Env{"FOO": "BAR"}}},
			},
		}
		state, err := with.Compile(llb.Scratch(), ChainStates{})
		req.NoError(err)

		foo, ok, err := state.GetEnv(context.TODO(), "FOO")
		req.NoError(err)
		req.True(ok)
		req.Equal("BAR", foo)
	})

	t.Run("WorkingDirectory", func(t *testing.T) {
		with := &With{
			With: Options{
				{WorkingDirectory: &WorkingDirectory{Directory: "/src"}},
			},
		}
		state, err := with.Compile(llb.Scratch(), ChainStates{})
		req.NoError(err)

		dir, err := state.GetDir(context.TODO())
		req.NoError(err)
		req.Equal("/src", dir)

		state = state.Run(llb.Shlex("foo")).Root()

		def, err := state.Marshal(context.TODO())
		req.NoError(err)

		llbreq := llbtest.New(t, def)

		_, eops := llbreq.ContainsNExecOps(1)
		req.Equal("/src", eops[0].Exec.Meta.Cwd)
	})
}
