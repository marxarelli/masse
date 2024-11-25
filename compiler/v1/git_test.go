package v1

import (
	"testing"

	"gitlab.wikimedia.org/dduvall/masse/util/llbtest"
	"gitlab.wikimedia.org/dduvall/masse/util/testcompile"
)

func TestGit(t *testing.T) {
	compile := testcompile.New(
		t,
		[]string{"wikimedia.org/dduvall/masse/state"},
		testcompile.WithCompiler(func() *compiler {
			return newCompiler(nil)
		}),
	)

	compile.Test(
		"minimal",
		`state.#Git & { git: "https://an.example/repo.git" }`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			_, sops := req.ContainsNSourceOps(1)
			req.Equal("git://an.example/repo.git#refs/heads/main", sops[0].Source.Identifier)
			req.Contains(sops[0].Source.Attrs, "git.fullurl")
			req.Equal("https://an.example/repo.git", sops[0].Source.Attrs["git.fullurl"])
		},
	)

	compile.Test(
		"ref",
		`state.#Git & { git: "https://an.example/repo.git", ref: "refs/tags/foo" }`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			_, sops := req.ContainsNSourceOps(1)
			req.Equal("git://an.example/repo.git#refs/tags/foo", sops[0].Source.Identifier)
			req.Contains(sops[0].Source.Attrs, "git.fullurl")
			req.Equal("https://an.example/repo.git", sops[0].Source.Attrs["git.fullurl"])
		},
	)

	compile.Test(
		"options/keepGitdir/true",
		`state.#Git & { git: "https://an.example/repo.git", options: keepGitDir: true }`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			_, sops := req.ContainsNSourceOps(1)
			req.Equal("git://an.example/repo.git#refs/heads/main", sops[0].Source.Identifier)
			req.Contains(sops[0].Source.Attrs, "git.keepgitdir")
			req.Equal("true", sops[0].Source.Attrs["git.keepgitdir"])
		},
	)

	compile.Test(
		"options/keepGitdir/false",
		`state.#Git & { git: "https://an.example/repo.git", options: keepGitDir: false }`,
		func(t *testing.T, req *llbtest.Assertions, _ *testcompile.Test) {
			_, sops := req.ContainsNSourceOps(1)
			req.Equal("git://an.example/repo.git#refs/heads/main", sops[0].Source.Identifier)
			req.NotContains(sops[0].Source.Attrs, "git.keepgitdir")
		},
	)
}
