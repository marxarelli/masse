package target

import (
	"wikimedia.org/dduvall/phyton/schema/common"
	"wikimedia.org/dduvall/phyton/schema/state"
)

#Target: {
	build!: state.#ChainRef
	platforms!: [#TargetPlatform, ...#TargetPlatform]
	platformsValue: [
		for p in platforms {
			[
				if (p & string) != _|_ {common.#Platform & {name: p}},
				p,
			][0]
		}
	]
	labels: common.#Labels
	runtime: #Runtime
}

#TargetPlatform: string | common.#Platform
