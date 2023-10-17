package state

import (
	"wikimedia.org/dduvall/phyton/schema/common"
)

#RunOption: {
	#Host |
	#CacheMount | #SourceMount |
	#TmpFSMount | #ReadOnly |
	#Option
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
	from:   #ChainRef
	source: string | *"/"
}

#TmpFSMount: {
	tmpfs!: string
	size:   uint64
}

#ReadOnly: {
	readOnly!: *true | false
}
