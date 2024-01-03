package state

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.wikimedia.org/dduvall/masse/util/llbtest"
	"gitlab.wikimedia.org/dduvall/masse/util/testdecode"
)

func TestDecodeGit(t *testing.T) {
	tester := &testdecode.Tester{
		T:          t,
		CUEImports: []string{"wikimedia.org/dduvall/masse/schema/state"},
	}

	testdecode.Run(tester,
		"state.#Git",
		`state.#Git & { git: "https://some.example/repo.git", ref: "refs/foo" }`,
		Git{
			Repo: "https://some.example/repo.git",
			Ref:  "refs/foo",
		},
	)

	testdecode.Run(tester,
		"state.#Git/options/keepGitDir",
		`state.#Git & { git: "https://some.example/repo.git", options: [ { keepGitDir: true } ] }`,
		Git{
			Repo: "https://some.example/repo.git",
			Ref:  "refs/heads/main",
			Options: []*GitOption{
				{KeepGitDir: &KeepGitDir{KeepGitDir: true}},
			},
		},
	)
}

func TestCompileGit(t *testing.T) {
	req := require.New(t)

	git := Git{
		Repo: "some.example/repo.git",
		Ref:  "refs/foo",
		Options: []*GitOption{
			{KeepGitDir: &KeepGitDir{KeepGitDir: true}},
		},
	}

	state, err := git.CompileSource(ChainStates{})
	req.NoError(err)

	def, err := state.Marshal(context.TODO())
	req.NoError(err)

	llbreq := llbtest.New(t, def)

	_, sops := llbreq.ContainsNSourceOps(1)
	req.Equal("git://some.example/repo.git#refs/foo", sops[0].Source.Identifier)
	req.Contains(sops[0].Source.Attrs, "git.keepgitdir")
	req.Equal("true", sops[0].Source.Attrs["git.keepgitdir"])
}
