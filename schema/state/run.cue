package state

import (
	"strings"
	"github.com/marxarelli/masse/common"
)

#Run: {
	run!:       [string, ...string]
	options?:   #RunOption | [#RunOption, ...#RunOption]

	#defaultOptions: [
		{ customName: strings.TrimSpace("💻 " + strings.Join(run, " ")) },
	]
}

#RunOption: {
	#Host |
	#CacheMount | #SourceMount | #TmpFSMount |
	#ValidExitCodes |
	#Option |
	#Constraint
}

#Host: {
	common.#Host
}

#CacheMount: {
	cache!: string
	access: *"shared" | "private" | "locked"
}

#SourceMount: {
	mount!:   string
	from:     #ChainRef
	source:   string | *"/"
	readonly: bool | *false
}

#TmpFSMount: {
	tmpfs!: string
	size:   uint64 | *100Mi
}

#ValidExitCodes: {
	validExitCodes!: [uint64, ...uint64]
}
