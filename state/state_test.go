package state

import (
	"testing"

	"gitlab.wikimedia.org/dduvall/phyton/common"
	"gitlab.wikimedia.org/dduvall/phyton/util/testdecode"
)

func TestDecodeState(t *testing.T) {
	tester := &testdecode.Tester{
		T:          t,
		CUEImports: []string{"wikimedia.org/dduvall/phyton/schema/state"},
	}

	testdecode.Run(tester,
		"state.#Chain",
		`state.#Chain & [
			{ image: "foo.example/image/ref" },
			{ with: [ { env: { FOO: "BAR" } } ] },
			{ run: "apt-get install foo", options: [ { cache: "/var/cache/apt", access: "locked" } ] },
			{ copy: ["foo/*"], from: "foo", options: [ { followSymlinks: true } ] },
			{ diff: [ { run: "make stuff" } ] },
		]`,
		Chain{
			{Image: &Image{Ref: "foo.example/image/ref"}},
			{With: &With{With: []*Option{
				{Env: &Env{Env: common.Env{"FOO": "BAR"}}},
			}}},
			{Run: &Run{
				Command: "apt-get install foo",
				Options: []*RunOption{
					{CacheMount: &CacheMount{Target: "/var/cache/apt", Access: CacheLocked}},
				},
			}},
			{Copy: &Copy{
				Source:      []common.Glob{"foo/*"},
				Destination: "./",
				From:        "foo",
				Options: []*CopyOption{
					{FollowSymlinks: &FollowSymlinks{FollowSymlinks: true}},
				},
			}},
			{Diff: &Diff{
				Upper: Chain{
					{Run: &Run{Command: "make stuff"}},
				},
			}},
		},
	)
}