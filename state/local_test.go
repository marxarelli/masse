package state

import (
	"testing"

	"gitlab.wikimedia.org/dduvall/phyton/common"
	"gitlab.wikimedia.org/dduvall/phyton/util/testdecode"
)

func TestDecodeLocal(t *testing.T) {
	tester := &testdecode.Tester{
		T:          t,
		CUEImports: []string{"wikimedia.org/dduvall/phyton/schema/state"},
	}

	testdecode.Run(tester,
		"state.#Local",
		`state.#Local & { local: "context" }`,
		Local{Name: "context"},
	)

	testdecode.Run(tester,
		"state.#Local/options/differ",
		`state.#Local & { local: "context", options: [ { differ: "metadata", require: true } ] }`,
		Local{
			Name: "context",
			Options: []*LocalOption{
				{Differ: &Differ{Differ: "metadata", Require: true}},
			},
		},
	)

	testdecode.Run(tester,
		"state.#Local/options/exclude",
		`state.#Local & { local: "context", options: [ { exclude: ["/src/foo"] } ] }`,
		Local{
			Name: "context",
			Options: []*LocalOption{
				{Exclude: &Exclude{Exclude: []common.Glob{"/src/foo"}}},
			},
		},
	)

	testdecode.Run(tester,
		"state.#Local/options/followPaths",
		`state.#Local & { local: "context", options: [ { followPaths: ["/src/foo"] } ] }`,
		Local{
			Name: "context",
			Options: []*LocalOption{
				{FollowPaths: &FollowPaths{FollowPaths: []string{"/src/foo"}}},
			},
		},
	)

	testdecode.Run(tester,
		"state.#Local/options/include",
		`state.#Local & { local: "context", options: [ { include: ["/src/foo"] } ] }`,
		Local{
			Name: "context",
			Options: []*LocalOption{
				{Include: &Include{Include: []common.Glob{"/src/foo"}}},
			},
		},
	)

	testdecode.Run(tester,
		"state.#Local/options/sharedKeyHint",
		`state.#Local & { local: "context", options: [ { sharedKeyHint: "foo-hint" } ] }`,
		Local{
			Name: "context",
			Options: []*LocalOption{
				{SharedKeyHint: &SharedKeyHint{SharedKeyHint: "foo-hint"}},
			},
		},
	)
}
