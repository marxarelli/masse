package state

import (
	"wikimedia.org/dduvall/masse/common"
)

// Local uses a directory from the build host as the initial filesystem for a build chain.
#Local: {
	local!:   string
	options?: #LocalOption | [#LocalOption, ...#LocalOption]

	#defaultOptions: [
		{ customName: "💻 \(local)" },
	]
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
