package state

import (
	"wikimedia.org/dduvall/masse/common"
)

#Copy: {
	copy!:       string
	from!:       #ChainRef
	destination: string | *"./"
	options?:    #CopyOption | [#CopyOption, ...#CopyOption]
}

#CopyOption: {
	common.#Creation |
	common.#User | common.#Group | common.#Mode |
	common.#Include | common.#Exclude |
	#FollowSymlinks | #CopyDirectoryContents |
	#AllowNotFound | #Wildcard |
	#ReplaceExisting |
	#CreateParents
}

#FollowSymlinks: {
	followSymlinks!: *true | false
}

#CopyDirectoryContents: {
	copyDirectoryContents!: true | *false
}

#ReplaceExisting: {
	replaceExisting!: true | *false
}
