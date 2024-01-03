package state

import (
	"list"
	"wikimedia.org/dduvall/masse/schema/common"
)

#Local: {
	local!: string
	options?: #LocalOption | [#LocalOption, ...#LocalOption]
	if options != _|_ {
		optionsValue: list.FlattenN([options], 1)
	}
}

#LocalOption: {
	common.#Include | common.#Exclude |
	#FollowPaths | #SharedKeyHint |
	#Differ |
	#Constraint
}

#FollowPaths: {
	followPaths!: [...string]
}

#SharedKeyHint: {
	sharedKeyHint!: string
}

#Differ: {
	differ!: *"metadata" | "none"
	require: *true | false
}
