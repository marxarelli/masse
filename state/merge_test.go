package state

import (
	"testing"

	"gitlab.wikimedia.org/dduvall/phyton/util/testdecode"
)

func TestDecodeMerge(t *testing.T) {
	tester := &testdecode.Tester{
		T:          t,
		CUEImports: []string{"wikimedia.org/dduvall/phyton/schema/state"},
	}

	testdecode.Run(tester,
		"state.#Merge",
		`state.#Merge & { merge: ["foo", "bar"] }`,
		Merge{
			Merge: []ChainRef{"foo", "bar"},
		},
	)
}
