package state

import (
	"list"
)

#CopyOp: {
	copy!:       string | [string, ...string]
	from!:       #ChainRef
	destination: string | *"./"
	options?:    #CopyOption | [#CopyOption, ...#CopyOption]

	_from: from
	_destination: destination
	_options: options

	ops: [
		{
			file: [
				for cp in list.FlattenN([copy], 1) {
					copy: cp
					from: _from
					destination: _destination
					if _options != _|_ {
						options: _options
					}
				}
			]
		}
	]
}
