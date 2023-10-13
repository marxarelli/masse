package state

import (
	"testing"

	"gitlab.wikimedia.org/dduvall/phyton/util/testdecode"
)

func TestDecodeScratch(t *testing.T) {
	tester := &testdecode.Tester{
		T:          t,
		CUEImports: []string{"wikimedia.org/dduvall/phyton/schema/state"},
	}

	testdecode.Run(tester,
		"state.#Scratch",
		`state.#Scratch & { scratch: true }`,
		Scratch{Scratch: true},
	)
}
