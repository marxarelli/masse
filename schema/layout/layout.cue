package layout

import (
	"wikimedia.org/dduvall/phyton/schema/common"
	"wikimedia.org/dduvall/phyton/schema/state"
)

#Layout: {
	comprises!: [state.#ChainRef, ...state.#ChainRef]
	authors: [#Author, ...#Author]
	platforms?: [common.#Platform, ...common.#Platform]
	parameters?:    #Parameters
	configuration?: #ImageConfig
}
