package state

import (
	"wikimedia.org/dduvall/phyton/common"
)

#CopyOption: {
	common.#Creation |
	common.#User | common.#Group | common.#Mode |
	common.#Include | common.#Exclude |
	#FollowSymlinks | #CopyDirectoryContent
}

#FollowSymlinks: {
	follow_symlinks!: *true | false
}

#CopyDirectoryContent: {
	copy_directory_content!: true | *false
}
