package state

import (
	"wikimedia.org/dduvall/phyton/common"
)

#RunOption: {
	#Env | #Host |
	#CacheMount | #SourceMount |
	#TmpFsMount | #ReadonlyRoot
}

#Host: {
	common.#Host
}

#CacheMount: {
	cache!: string
	access: *"shared" | "private" | "locked"
}

#SourceMount: {
	mount!: string
	from:   #Chain | *null
	source: string | *"/"
}

#TmpFsMount: {
	tmpfs!: string
	size:   uint64
}

#ReadonlyRoot: {
	readonly!: *true | false
}
