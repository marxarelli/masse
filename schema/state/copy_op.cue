package state

import (
	"list"
)

#CopyOp: {
	copy!:                      string | [string, ...string]
	$from        = from!:       #ChainRef
	$destination = destination: string | *"./"
	$options     = options?:    #CopyOption | [#CopyOption, ...#CopyOption]

	let sources = list.FlattenN([copy], 1)

	ops: [
		{
			file: [
				for src in sources {
					copy: src
					from: $from
					destination: $destination
					options?: $options
				}
			]
		}
	]
}
