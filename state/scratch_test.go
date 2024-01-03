package state

import (
	"testing"

	"gitlab.wikimedia.org/dduvall/masse/util/testdecode"
)

func TestDecodeScratch(t *testing.T) {
	tester := &testdecode.Tester{
		T:          t,
		CUEImports: []string{"wikimedia.org/dduvall/masse/schema/state"},
	}

	testdecode.Run(tester,
		"state.#Scratch",
		`state.#Scratch & { scratch: true }`,
		Scratch{Scratch: true},
	)
}
