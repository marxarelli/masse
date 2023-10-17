package state

import (
	"list"
)

#Link: {
	link!:       string | [string, ...string]
	source:      list.FlattenN([link], 1)
	from:        #ChainRef
	destination: string | *"./"

	options?: #CopyOption | [#CopyOption, ...#CopyOption]
	if options != _|_ {
		optionsValue: list.FlattenN([options], 1)
	}
}
