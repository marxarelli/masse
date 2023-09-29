package state

import (
	"wikimedia.org/dduvall/phyton/common"
)

#Image: {
	image!: string

	options?: [...#ImageOption]
}

#ImageOption: {
	common.#Platform | #LayerLimit
}

#LayerLimit: {
	layer_limit: uint32
}
