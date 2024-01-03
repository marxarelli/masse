package common

import (
	"net"
	"testing"
	"time"

	"gitlab.wikimedia.org/dduvall/masse/util/testdecode"
)

func TestDecode(t *testing.T) {
	tester := &testdecode.Tester{
		T:          t,
		CUEImports: []string{"wikimedia.org/dduvall/masse/schema/common"},
	}

	ctime, _ := time.Parse(time.RFC3339, "2020-01-20T01:02:03Z")
	testdecode.Run(tester,
		"common.#Creation",
		`common.#Creation & { ctime: "2020-01-20T01:02:03Z" }`,
		Creation{Ctime: &ctime},
	)

	testdecode.Run(tester,
		"common.#Env",
		`common.#Env & { foo: "BAR" }`,
		Env{"foo": "BAR"},
	)

	testdecode.Run(tester,
		"common.#Exclude",
		`common.#Exclude & { exclude: ["foo/*.bar"] }`,
		Exclude{Exclude: []Glob{"foo/*.bar"}},
	)

	gid := uint32(123)
	testdecode.Run(tester,
		"common.#Group/gid",
		`common.#Group & { gid: 123 }`,
		Group{GID: &gid, Group: ""},
	)

	testdecode.Run(tester,
		"common.#Group/group",
		`common.#Group & { group: "foo" }`,
		Group{GID: nil, Group: "foo"},
	)

	testdecode.Run(tester,
		"common.#Host",
		`common.#Host & { ip: "1.2.3.4", host: "foo.example" }`,
		Host{IP: net.ParseIP("1.2.3.4"), Host: "foo.example"},
	)

	testdecode.Run(tester,
		"common.#Include",
		`common.#Include & { include: ["foo/*.bar"] }`,
		Include{Include: []Glob{"foo/*.bar"}},
	)

	testdecode.Run(tester,
		"common.#Labels",
		`common.#Labels & { "foo.label": "bar" }`,
		Labels{"foo.label": "bar"},
	)

	testdecode.Run(tester,
		"common.#Mode/numeric",
		`common.#Mode & { mode: 0o2755 }`,
		Mode{Mode: 0o2755},
	)

	testdecode.Run(tester,
		"common.#Mode/symbolic",
		`common.#Mode & { mode: "rwxr-sr-x" }`,
		Mode{Mode: 0o2755},
	)

	testdecode.Run(tester,
		"common.#Platform",
		`common.#Platform & { os: "linux", architecture: "arm64", variant: "v8" }`,
		Platform{OS: "linux", Architecture: "arm64", Variant: "v8"},
	)

	uid := uint32(123)
	testdecode.Run(tester,
		"common.#User/uid",
		`common.#User & { uid: 123 }`,
		User{UID: &uid, User: ""},
	)

	testdecode.Run(tester,
		"common.#User/user",
		`common.#User & { user: "foo" }`,
		User{UID: nil, User: "foo"},
	)
}
