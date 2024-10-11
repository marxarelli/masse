package target

import (
	"wikimedia.org/dduvall/masse/common"
	"wikimedia.org/dduvall/masse/state"
)

#Target: {
	build!: state.#Chain
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
