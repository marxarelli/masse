package state

import (
	"github.com/marxarelli/masse/common"
)

#Mkfile: {
	mkfile!:  string
	content!: bytes
	options?: #MkfileOption | [#MkfileOption, ...#MkfileOption]
}

#MkfileOption: {
	common.#Creation |
	common.#User | common.#Group |
	common.#Mode
}
