package state

import (
	"list"
)

#Copy: {
	copy!:       string | [string, ...string]
	source:      list.FlattenN([copy], 1)
	from:        #ChainRef
	destination: string | *"./"

	options?: #CopyOption | [#CopyOption, ...#CopyOption]
	if options != _|_ {
		optionsValue: list.FlattenN([options], 1)
	}
}
