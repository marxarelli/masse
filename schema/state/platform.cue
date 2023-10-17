package state

import (
	"strings"
	"wikimedia.org/dduvall/phyton/schema/common"
)

#Platform: {
	#SymbolicPlatform | #LiteralPlatform
}

#SymbolicPlatform: {
	platform!: string
	_parts:    strings.SplitN(platform, "/", 3)
	platformValue:     common.#Platform & [
		if len(_parts) > 2 {
			{os: _parts[0], architecture: _parts[1], variant: _parts[2]}
		},
		{os: _parts[0], architecture: _parts[1]},
	][0]
}

#LiteralPlatform: {
	platform!: common.#Platform
	platformValue:     platform
}
