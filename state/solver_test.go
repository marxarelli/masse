package state

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.wikimedia.org/dduvall/phyton/common"
)

func TestSolve(t *testing.T) {
	req := require.New(t)

	chains := Chains{
		"repo": Chain{
			&State{Git: &Git{Repo: "some.example/repo.git"}},
		},
		"build": Chain{
			&State{Image: &Image{Ref: "some.example/build/env"}},
		},
		"tools": Chain{
			&State{Merge: &Merge{Merge: []ChainRef{"build"}}},
			&State{Diff: &Diff{Upper: Chain{
				&State{Run: &Run{Command: "apt-get install build-essential"}},
			}}},
		},
		"binaries": Chain{
			&State{Merge: &Merge{Merge: []ChainRef{"build", "tools"}}},
			&State{With: &With{With: []*Option{
				{WorkingDirectory: &WorkingDirectory{Directory: "/src"}},
			}}},
			&State{Run: &Run{
				Command: "make foo && cp foo /srv/foo",
				Options: []*RunOption{
					{SourceMount: &SourceMount{Target: "/src", From: "repo"}},
				},
			}},
		},
		"final": Chain{
			&State{Copy: &Copy{Source: []common.Glob{"/srv/foo"}, From: "binaries"}},
		},
		"production": Chain{
			&State{Image: &Image{Ref: "some.example/prod/env"}},
		},
	}

	g, err := NewGraph(chains, &Merge{Merge: []ChainRef{"production", "final"}})
	req.NoError(err)

	state, err := NewSolver().Solve(g)
	req.NoError(err)
	req.NotNil(state)
}
