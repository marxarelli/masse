package state

import (
	"wikimedia.org/dduvall/masse/common"
)

#Local: {
	local!:   string
	options?: #LocalOption | [#LocalOption, ...#LocalOption]
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
