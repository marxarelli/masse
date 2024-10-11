package state

import (
	"list"
)

#Copy: {
	copy!:       string | [string, ...string]
	source:      list.FlattenN([copy], 1)
	from:        #ChainRef
	destination: string | *"./"
	options?:    #CopyOption | [#CopyOption, ...#CopyOption]
}
