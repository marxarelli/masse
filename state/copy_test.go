package state

import (
	"testing"
	"time"

	"gitlab.wikimedia.org/dduvall/phyton/common"
	"gitlab.wikimedia.org/dduvall/phyton/util/testdecode"
)

func TestDecodeCopy(t *testing.T) {
	tester := &testdecode.Tester{
		T:          t,
		CUEImports: []string{"wikimedia.org/dduvall/phyton/schema/state"},
	}

	testdecode.Run(tester,
		"state.#Copy",
		`state.#Copy & { copy: "foo/*", from: "local" }`,
		Copy{
			Source:      []common.Glob{"foo/*"},
			Destination: "./",
			From:        "local",
		},
	)

	testdecode.Run(tester,
		"state.#Copy/destination",
		`state.#Copy & { copy: "foo/*", from: "local", destination: "/srv/foo" }`,
		Copy{
			Source:      []common.Glob{"foo/*"},
			Destination: "/srv/foo",
			From:        "local",
		},
	)

	ctime, _ := time.Parse(time.RFC3339, "2020-01-20T01:02:03Z")
	testdecode.Run(tester,
		"state.#Copy/options/creation",
		`state.#Copy & { copy: ["foo/*"], from: "foo", options: [ { ctime: "2020-01-20T01:02:03Z" } ] }`,
		Copy{
			Source:      []common.Glob{"foo/*"},
			Destination: "./",
			From:        "foo",
			Options: []*CopyOption{
				{Creation: &common.Creation{Ctime: &ctime}},
			},
		},
	)

	testdecode.Run(tester,
		"state.#Copy/options/user",
		`state.#Copy & { copy: ["foo/*"], from: "foo", options: [ { user: "foo" } ] }`,
		Copy{
			Source:      []common.Glob{"foo/*"},
			Destination: "./",
			From:        "foo",
			Options: []*CopyOption{
				{User: &common.User{User: "foo"}},
			},
		},
	)

	testdecode.Run(tester,
		"state.#Copy/options/group",
		`state.#Copy & { copy: ["foo/*"], from: "foo", options: [ { group: "foo" } ] }`,
		Copy{
			Source:      []common.Glob{"foo/*"},
			Destination: "./",
			From:        "foo",
			Options: []*CopyOption{
				{Group: &common.Group{Group: "foo"}},
			},
		},
	)

	testdecode.Run(tester,
		"state.#Copy/options/mode",
		`state.#Copy & { copy: ["foo/*"], from: "foo", options: [ { mode: "rwxr-xr-x" } ] }`,
		Copy{
			Source:      []common.Glob{"foo/*"},
			Destination: "./",
			From:        "foo",
			Options: []*CopyOption{
				{Mode: &common.Mode{Mode: 0o0755}},
			},
		},
	)

	testdecode.Run(tester,
		"state.#Copy/options/include",
		`state.#Copy & { copy: ["foo/*"], from: "foo", options: [ { include: ["*.sh"] } ] }`,
		Copy{
			Source:      []common.Glob{"foo/*"},
			Destination: "./",
			From:        "foo",
			Options: []*CopyOption{
				{Include: &common.Include{Include: []common.Glob{"*.sh"}}},
			},
		},
	)

	testdecode.Run(tester,
		"state.#Copy/options/exclude",
		`state.#Copy & { copy: ["foo/*"], from: "foo", options: [ { exclude: ["*.sh"] } ] }`,
		Copy{
			Source:      []common.Glob{"foo/*"},
			Destination: "./",
			From:        "foo",
			Options: []*CopyOption{
				{Exclude: &common.Exclude{Exclude: []common.Glob{"*.sh"}}},
			},
		},
	)

	testdecode.Run(tester,
		"state.#Copy/options/followSymlinks",
		`state.#Copy & { copy: ["foo/*"], from: "foo", options: [ { followSymlinks: true } ] }`,
		Copy{
			Source:      []common.Glob{"foo/*"},
			Destination: "./",
			From:        "foo",
			Options: []*CopyOption{
				{FollowSymlinks: &FollowSymlinks{FollowSymlinks: true}},
			},
		},
	)

	testdecode.Run(tester,
		"state.#Copy/options/copyDirectoryContent",
		`state.#Copy & { copy: ["foo/*"], from: "foo", options: [ { copyDirectoryContent: true } ] }`,
		Copy{
			Source:      []common.Glob{"foo/*"},
			Destination: "./",
			From:        "foo",
			Options: []*CopyOption{
				{CopyDirectoryContent: &CopyDirectoryContent{CopyDirectoryContent: true}},
			},
		},
	)
}
