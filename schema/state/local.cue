package state

import (
	"wikimedia.org/dduvall/phyton/schema/common"
)

#Local: {
	local!: string
	options?: [...#LocalOption]
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
