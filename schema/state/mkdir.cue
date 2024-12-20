package state

import (
	"github.com/marxarelli/masse/common"
)

#Mkdir: {
	mkdir!:   string
	options?: #MkdirOption | [#MkdirOption, ...#MkdirOption]
}

#MkdirOption: {
	#MkfileOption | #CreateParents |
	common.#Mode
}

#CreateParents: {
	createParents!: bool
}
