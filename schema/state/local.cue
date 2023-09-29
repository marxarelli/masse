package state

import (
	"wikimedia.org/dduvall/phyton/common"
)

#Local: {
	local!: string
	options?: [...#LocalOption]
}

#LocalOption: {
	common.#Include | common.#Exclude | #FollowPaths | #SharedKeyHint | #Differ
}

#FollowPaths: {
	follow_paths: [...string]
}

#SharedKeyHint: {
	shared_key_hint: string
}

#Differ: {
	type:    *"metadata" | "none"
	require: *true | false
}
