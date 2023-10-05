package state

import (
	"testing"

	"gitlab.wikimedia.org/dduvall/phyton/util/testdecode"
)

func TestDecodeGit(t *testing.T) {
	tester := &testdecode.Tester{
		T:          t,
		CUEImports: []string{"wikimedia.org/dduvall/phyton/schema/state"},
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
