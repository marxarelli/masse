package state

import (
	"net"
	"testing"

	"gitlab.wikimedia.org/dduvall/masse/common"
	"gitlab.wikimedia.org/dduvall/masse/util/testdecode"
)

func TestDecode(t *testing.T) {
	tester := &testdecode.Tester{
		T:          t,
		CUEImports: []string{"wikimedia.org/dduvall/masse/schema/state"},
	}

	testdecode.Run(tester,
		"state.#CacheMount",
		`state.#CacheMount & { cache: "/var/cache/foo", access: "locked" }`,
		CacheMount{Target: "/var/cache/foo", Access: CacheLocked},
	)

	testdecode.Run(tester,
		"state.#Env",
		`state.#Env & { env: { "foo": "BAR" } }`,
		Env{Env: common.Env{"foo": "BAR"}},
	)

	testdecode.Run(tester,
		"state.#Host",
		`state.#Host & { host: "foo.example", ip: "1.2.3.4" }`,
		Host{Host: "foo.example", IP: net.ParseIP("1.2.3.4")},
	)

	testdecode.Run(tester,
		"state.#LayerLimit",
		`state.#LayerLimit & { layerLimit: 999 }`,
		LayerLimit{LayerLimit: uint32(999)},
	)

	testdecode.Run(tester,
		"state.#Option/env",
		`state.#Option & { env: { "foo": "BAR" } }`,
		Option{Env: &Env{Env: common.Env{"foo": "BAR"}}},
	)

	testdecode.Run(tester,
		"state.#Option/directory",
		`state.#Option & { directory: "/srv/foo" }`,
		Option{WorkingDirectory: &WorkingDirectory{Directory: "/srv/foo"}},
	)

	testdecode.Run(tester,
		"state.#ReadOnly",
		`state.#ReadOnly & { readOnly: true }`,
		ReadOnly{ReadOnly: true},
	)

	testdecode.Run(tester,
		"state.#TmpFSMount",
		`state.#TmpFSMount & { tmpfs: "/var/cache/foo", size: 500Mi }`,
		TmpFSMount{TmpFS: "/var/cache/foo", Size: 500 * 1024 * 1024},
	)

	testdecode.Run(tester,
		"state.#WorkingDirectory",
		`state.#WorkingDirectory & { directory: "/srv/foo" }`,
		WorkingDirectory{Directory: "/srv/foo"},
	)
}
