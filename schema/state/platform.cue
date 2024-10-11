package state

import (
	"wikimedia.org/dduvall/masse/common"
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
