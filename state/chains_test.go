package state

import (
	"testing"

	"gitlab.wikimedia.org/dduvall/phyton/util/testdecode"
)

func TestDecodeChains(t *testing.T) {
	tester := &testdecode.Tester{
		T: t,
		CUEImports: []string{
			"wikimedia.org/dduvall/phyton/schema/state",
		},
	}

	testdecode.Run(tester,
		"state.#Chains",
		`state.#Chains & {
			foo: [
				{ git: "some.example/repo.git", ref: "foo" },
			]
			bar: [
				{ merge: ["foo"] },
				{ run: "foo" },
			]
			baz: [
				{ scratch: true },
			]
		}`,
		Chains{
			"foo": Chain{
				{Git: &Git{Repo: "some.example/repo.git", Ref: "foo"}},
			},
			"bar": Chain{
				{Merge: &Merge{Merge: []ChainRef{"foo"}}},
				{Run: &Run{Command: "foo"}},
			},
			"baz": Chain{
				{Scratch: &Scratch{Scratch: true}},
			},
		},
	)
}
