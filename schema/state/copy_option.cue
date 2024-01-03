package state

import (
	"wikimedia.org/dduvall/masse/schema/common"
)

#CopyOption: {
	common.#Creation |
	common.#User | common.#Group | common.#Mode |
	common.#Include | common.#Exclude |
	#FollowSymlinks | #CopyDirectoryContent
}

#FollowSymlinks: {
	followSymlinks!: *true | false
}

#CopyDirectoryContent: {
	copyDirectoryContent!: true | *false
}
