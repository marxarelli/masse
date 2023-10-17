package state

import (
	"list"
)

#With: {
	with!: #Option | [#Option, ...#Option]
	withValue: list.FlattenN([with], 1)
}
