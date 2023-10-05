package layout

import (
	"testing"

	"gitlab.wikimedia.org/dduvall/phyton/state"
	"gitlab.wikimedia.org/dduvall/phyton/util/testdecode"
)

func TestDecodeChains(t *testing.T) {
	tester := &testdecode.Tester{
		T: t,
		CUEImports: []string{
			"wikimedia.org/dduvall/phyton/schema/layout",
		},
	}

	testdecode.Run(tester,
		"layout.#Chains",
		`layout.#Chains & {
			foo: [
				{ git: "some.example/repo.git", ref: "foo" },
			]
			bar: [
				{ merge: ["foo"] },
				{ run: "foo" },
			]
		}`,
		Chains{
			"foo": state.Chain{
				{Git: &state.Git{Repo: "some.example/repo.git", Ref: "foo"}},
			},
			"bar": state.Chain{
				{Merge: &state.Merge{Merge: []state.ChainRef{"foo"}}},
				{Run: &state.Run{Command: "foo"}},
			},
		},
	)
}
