package common

import (
	"net"
	"testing"
	"time"

	"cuelang.org/go/cue"
	"github.com/stretchr/testify/require"
	"gitlab.wikimedia.org/dduvall/masse/util/testconfig"
)

func testDecodeCase[T any](tester *testconfig.Tester, name string, expr string, expected T) {
	tester.Test(
		name,
		expr,
		func(t *testing.T, req *require.Assertions, v cue.Value) {
			t.Helper()
			actual := *new(T)
			req.NoError(v.Decode(&actual))
			req.Equal(actual, expected)
		},
	)
}

func TestDecode(t *testing.T) {
	tester := testconfig.New(t, []string{"wikimedia.org/dduvall/masse/common"})

	ctime, _ := time.Parse(time.RFC3339, "2020-01-20T01:02:03Z")
	testDecodeCase(tester,
		"common.#Creation",
		`common.#Creation & { ctime: "2020-01-20T01:02:03Z" }`,
		Creation{Ctime: &ctime},
	)

	testDecodeCase(tester,
		"common.#Env",
		`common.#Env & { foo: "BAR" }`,
		Env{"foo": "BAR"},
	)

	testDecodeCase(tester,
		"common.#Exclude",
		`common.#Exclude & { exclude: ["foo/*.bar"] }`,
		Exclude{Exclude: []Glob{"foo/*.bar"}},
	)

	gid := uint32(123)
	testDecodeCase(tester,
		"common.#Group/gid",
		`common.#Group & { gid: 123 }`,
		Group{GID: &gid, Group: ""},
	)

	testDecodeCase(tester,
		"common.#Group/group",
		`common.#Group & { group: "foo" }`,
		Group{GID: nil, Group: "foo"},
	)

	testDecodeCase(tester,
		"common.#Host",
		`common.#Host & { ip: "1.2.3.4", host: "foo.example" }`,
		Host{IP: net.ParseIP("1.2.3.4"), Host: "foo.example"},
	)

	testDecodeCase(tester,
		"common.#Include",
		`common.#Include & { include: ["foo/*.bar"] }`,
		Include{Include: []Glob{"foo/*.bar"}},
	)

	testDecodeCase(tester,
		"common.#Labels",
		`common.#Labels & { "foo.label": "bar" }`,
		Labels{"foo.label": "bar"},
	)

	testDecodeCase(tester,
		"common.#Mode/numeric",
		`common.#Mode & { mode: 0o2755 }`,
		Mode{Mode: 0o2755},
	)

	testDecodeCase(tester,
		"common.#Mode/symbolic",
		`common.#Mode & { mode: "rwxr-sr-x" }`,
		Mode{Mode: 0o2755},
	)

	testDecodeCase(tester,
		"common.#Platform",
		`common.#Platform & { os: "linux", architecture: "arm64", variant: "v8" }`,
		Platform{OS: "linux", Architecture: "arm64", Variant: "v8"},
	)

	uid := uint32(123)
	testDecodeCase(tester,
		"common.#User/uid",
		`common.#User & { uid: 123 }`,
		User{UID: &uid, User: ""},
	)

	testDecodeCase(tester,
		"common.#User/user",
		`common.#User & { user: "foo" }`,
		User{UID: nil, User: "foo"},
	)
}
