package state

import (
	"github.com/marxarelli/masse/common"
)

#Platform: {
	#SymbolicPlatform | #LiteralPlatform
}

#SymbolicPlatform: {
	platform!:     string
	platformValue: common.#Platform & {name: platform}
}

#LiteralPlatform: {
	platform!:     common.#Platform
	platformValue: platform
}
