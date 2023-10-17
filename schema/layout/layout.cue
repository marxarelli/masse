package layout

import (
	"wikimedia.org/dduvall/phyton/schema/common"
	"wikimedia.org/dduvall/phyton/schema/state"
)

#Layout: {
	comprises!: [state.#ChainRef, ...state.#ChainRef]
	authors: [#Author, ...#Author]
	platforms?: [#LayoutPlatform, ...#LayoutPlatform]
	if platforms != _|_ {
		platformsValue: [
			for p in platforms {
				[
					if (p & string) != _|_ { common.#Platform & { name: p } },
					p,
				][0]
			}
		]
	}
	parameters?:    #Parameters
	configuration?: #ImageConfig
}

#LayoutPlatform: string | common.#Platform
