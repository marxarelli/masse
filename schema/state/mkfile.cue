package state

import (
	"wikimedia.org/releng/masse/common"
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
