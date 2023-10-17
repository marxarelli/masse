package state

import (
	"wikimedia.org/dduvall/phyton/schema/common"
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
